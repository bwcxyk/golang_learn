/*
@Author : YaoKun
@Time : 2023/7/28 16:42
*/

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

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

	// 设置连接池的最大连接数和空闲连接数
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)

	return db, nil
}

func readData(filePath string) ([]string, error) {

	// 打开 Excel 文件
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	// 读取 Sheet1 中的所有行
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	var keywords []string

	// 迭代每一行，读取关键字
	for i, row := range rows {
		// 跳过表头行
		if i == 0 {
			continue
		}

		// 假设关键字所在列为第一列
		keyword := row[0] // row[0]表示第一列
		keywords = append(keywords, keyword)
	}

	return keywords, nil
}

func queryDatabase(db *sql.DB, query string, keyword string, resultChan chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 执行查询
	// searchTerm 是查询参数，需要替换为你自己的查询参数
	rows, err := db.Query(query, sql.Named("searchTerm", keyword))
	if err != nil {
		log.Printf("查询关键字 %s 时发生错误: %v", keyword, err)
		resultChan <- []string{} // 发送空结果到channel，用于同步操作
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// 获取结果集的列信息
	columns, err := rows.Columns()
	if err != nil {
		log.Printf("获取列信息时发生错误: %v", err)
		resultChan <- []string{} // 发送空结果到channel，用于同步操作
		return
	}

	// 创建一个与列数相同长度的切片，用于存储每个字段的值
	values := make([]interface{}, len(columns))
	// 创建一个与列数相同长度的切片，用于存储每个字段的指针地址
	pointers := make([]interface{}, len(columns))

	// 为每个字段创建一个指针，并将其存储到pointers切片中
	for i := range values {
		pointers[i] = &values[i]
	}

	var result []string // 存储查询结果

	// 循环遍历结果集，并将每个字段的值读取到相应的指针地址
	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			log.Printf("扫描结果时发生错误: %v", err)
			resultChan <- []string{} // 发送空结果到channel，

			return
		}

		// 将每个字段的值转换为字符串，并添加到结果切片中
		for _, value := range values {
			result = append(result, fmt.Sprintf("%v", value))
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("迭代结果集时发生错误: %v", err)
		resultChan <- []string{} // 发送空结果到channel，用于同步操作
		return
	}

	// 打印查询结果
	//log.Printf("查询结果: %v", result)

	// 将查询结果发送到channel，用于同步操作
	resultChan <- result
}

func writeData(filePath string, results [][]string) error {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}

	sheetName := "Sheet2"
	_, err = file.NewSheet(sheetName)
	if err != nil {
		return err
	}

	for i, row := range results {
		for j, value := range row {
			colAlpha := convertToAlphaString(j + 1)
			cell := colAlpha + strconv.Itoa(i+1)
			err := file.SetCellValue(sheetName, cell, value)
			if err != nil {
				return err
			}
		}
	}

	err = file.SaveAs(filePath)
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

func main() {
	// 输入参数
	if len(os.Args) < 2 {
		log.Fatal("请提供 Excel 文件名作为命令行参数")
	}
	filePath := os.Args[1]

	// 打印日志
	log.Printf("连接数据库")

	// 连接数据库
	db, err := connectDb()
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// 打印日志
	log.Printf("读取关键字")

	// 读取关键字
	keywords, err := readData(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// 打印日志
	log.Printf("读取SQL")
	// 读取 SQL 文件
	sqlFile := "ora.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatal("无法读取 SQL 文件:", err)
	}
	sqlStatement := string(sqlBytes)

	// 创建结果切片，用于存储每个关键字的查询结果
	results := make([][]string, len(keywords))
	resultChan := make(chan []string) // 创建channel用于接收查询结果
	wg := sync.WaitGroup{}            // 创建WaitGroup用于等待所有goroutine执行完成

	// 打印日志
	log.Printf("查询数据")

	// 并发执行查询操作
	for i, keyword := range keywords {
		wg.Add(1)
		go queryDatabase(db, sqlStatement, keyword, resultChan, &wg)
		go func(index int) {
			results[index] = <-resultChan // 接收查询结果并存入对应索引位置
		}(i)
	}

	// 等待所有goroutine执行完成
	wg.Wait()

	// 打印日志
	log.Printf("写入 Excel 文件")

	// 将查询结果写回到原始 Excel 文件的 Sheet2
	err = writeData(filePath, results)
	if err != nil {
		log.Fatal(err)
	}
}
