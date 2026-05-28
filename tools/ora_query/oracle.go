/*
@Author : Yao Kun
@Time : 2023/7/28 16:42
*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	goora "github.com/sijms/go-ora/v2"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

// 定义一个结构体来匹配 YAML 文件的内容
type Config struct {
	Database struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		ServiceName string `yaml:"service_name"`
	} `yaml:"database"`
}

var (
	configOnce sync.Once
	config     Config
)

func initConfig() {
	// 读取 YAML 文件内容
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	// 解析 YAML 数据到 Config 结构体
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("error unmarshalling data: %v", err)
	}
}

func connectDb() (*sql.DB, error) {
	configOnce.Do(initConfig) // 确保配置只加载一次

	user := config.Database.User
	password := config.Database.Password
	host := config.Database.Host
	port := config.Database.Port
	serviceName := config.Database.ServiceName

	// 构建连接字符串
	dsn := goora.BuildUrl(host, port, serviceName, user, password, nil)

	// 使用连接字符串打开数据库连接
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func readData(filePath string) ([]string, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rows, err := file.GetRows("数据源")
	if err != nil {
		return nil, err
	}

	var keywords []string

	for i, row := range rows {
		if i == 0 || len(row) == 0 {
			continue
		}

		keyword := strings.TrimSpace(row[0])
		if keyword == "" {
			continue
		}
		keywords = append(keywords, keyword)
	}

	return keywords, nil
}

type queryResult struct {
	index   int
	columns []string
	result  []string
}

func queryDatabase(db *sql.DB, query string, keyword string, index int, resultChan chan<- queryResult) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, sql.Named("keyword", keyword))
	if err != nil {
		log.Printf("[错误] 查询关键字 %s 时发生错误: %v", keyword, err)
		resultChan <- queryResult{
			index:   index,
			columns: []string{},
			result:  []string{},
		}
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	columns, err := rows.Columns()
	if err != nil {
		log.Printf("[错误] 获取列信息时发生错误: %v", err)
		resultChan <- queryResult{
			index:   index,
			columns: []string{},
			result:  []string{},
		}
		return
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))

	for i := range values {
		pointers[i] = &values[i]
	}

	var result []string

	for rows.Next() {

		err := rows.Scan(pointers...)
		if err != nil {
			log.Printf("[错误] 扫描结果时发生错误: %v", err)
			resultChan <- queryResult{
				index:   index,
				columns: columns,
				result:  []string{},
			}
			return
		}

		for _, value := range values {
			result = append(result, valueToString(value))
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("[错误] 迭代结果集时发生错误: %v", err)
		resultChan <- queryResult{
			index:   index,
			columns: columns,
			result:  []string{},
		}
		return
	}

	resultChan <- queryResult{
		index:   index,
		columns: columns,
		result:  result,
	}
}

func writeData(filePath string, header []string, results [][]string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	sheetName := "查询_结果"
	_, err = file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	for j, value := range header {
		colAlpha := convertToAlphaString(j + 1)
		cell := colAlpha + "1"
		err := file.SetCellValue(sheetName, cell, value)
		if err != nil {
			return err
		}
	}

	for i, row := range results {
		for j, value := range row {
			colAlpha := convertToAlphaString(j + 1)
			cell := colAlpha + strconv.Itoa(i+2)
			err := file.SetCellValue(sheetName, cell, value)
			if err != nil {
				return err
			}
		}
	}

	err = file.Save()
	if err != nil {
		return err
	}

	return nil
}

func convertToAlphaString(n int) string {
	quotient := (n - 1) / 26
	remainder := (n - 1) % 26
	alpha := string('A' + rune(remainder))
	if quotient > 0 {
		alpha = convertToAlphaString(quotient) + alpha
	}
	return alpha
}

func valueToString(value interface{}) string {
	if value == nil {
		return ""
	}

	var strVal string

	switch v := value.(type) {
	case []byte:
		strVal = string(v)

	case time.Time:
		strVal = v.Format("2006-01-02 15:04:05")

	default:
		strVal = fmt.Sprintf("%v", v)
	}

	strVal = strings.TrimSpace(strVal)

	switch strings.ToLower(strVal) {
	case "", "<nil>", "null", "nil", "/*nil*/":
		return ""
	}

	return strVal
}

func normalizeSQL(raw string) string {
	s := strings.TrimSpace(raw)
	if strings.HasPrefix(s, "\uFEFF") {
		s = strings.TrimPrefix(s, "\uFEFF")
	}
	for strings.HasSuffix(s, ";") {
		s = strings.TrimSpace(strings.TrimSuffix(s, ";"))
	}
	if !utf8.ValidString(s) {
		return strings.ToValidUTF8(s, "")
	}
	return s
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("请提供 Excel 文件名作为命令行参数")
	}
	filePath := os.Args[1]

	log.Printf("连接数据库")
	db, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	log.Printf("读取关键字")
	keywords, err := readData(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("读取SQL")
	sqlFile := "search.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatal("无法读取 SQL 文件:", err)
	}
	sqlStatement := normalizeSQL(string(sqlBytes))
	if sqlStatement == "" {
		log.Fatal("SQL 文件内容为空")
	}

	results := make([][]string, len(keywords))
	headers := []string{}
	resultChan := make(chan queryResult, len(keywords))
	wg := sync.WaitGroup{}

	sem := make(chan struct{}, 5)

	log.Printf("开始查询数据")

	for i, keyword := range keywords {
		wg.Add(1)
		sem <- struct{}{} // 抢占名额

		go func(idx int, kw string) {
			defer wg.Done()
			defer func() { <-sem }() // 释放名额

			queryDatabase(db, sqlStatement, kw, idx, resultChan)
		}(i, keyword)
	}

	wg.Wait()
	close(resultChan)

	for item := range resultChan {
		results[item.index] = item.result
		if len(headers) == 0 && len(item.columns) > 0 {
			headers = item.columns
		}
	}

	log.Printf("写入 Excel 文件")
	err = writeData(filePath, headers, results)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("任务完成")
}