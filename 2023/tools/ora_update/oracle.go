/*
@Author : YaoKun
@Time : 2023/7/28 16:42
*/

package main

import (
	"database/sql"
	"log"
	"os"
	"sync"

	go_ora "github.com/sijms/go-ora/v2"
	"github.com/xuri/excelize/v2"
	"sigs.k8s.io/yaml"
)

// Config 结构体映射了YAML配置文件的结构
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Service  string `yaml:"service"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"database"`
}

// loadConfig 从指定的路径加载配置文件
func loadConfig(configPath string) (*Config, error) {
	var config Config

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// connectToDatabase 使用配置文件中的DSN来连接数据库
func connectToDatabase(configPath string) (*sql.DB, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}

	dsn := go_ora.BuildUrl(config.Database.Host, config.Database.Port, config.Database.Service, config.Database.Username, config.Database.Password, nil)

	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, err
	}
    // 设置连接池的最大连接数和空闲连接数
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)

	return db, nil
}

func readKeywordsFromExcel(filePath string) ([]string, error) {

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

		// 假设关键字所在列为 A 列
		keyword := row[0]
		keywords = append(keywords, keyword)
	}

	return keywords, nil
}

func updateDatabase(db *sql.DB, updateQuery string, keyword string, wg *sync.WaitGroup) {
	_, err := db.Exec(updateQuery, keyword)
	if err != nil {
		log.Printf("更新关键字 %s 时发生错误: %v", keyword, err)
		return
	}
	// 在成功执行后递减计数器
	log.Printf("更新关键字 %s 成功", keyword)
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
	db, err := connectToDatabase("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println("关闭数据库连接时发生错误:", err)
		}
	}(db)

	// 打印日志
	log.Printf("读取关键字")

	// 读取关键字
	keywords, err := readKeywordsFromExcel(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// 打印日志
	log.Printf("读取SQL")
	// 读取 SQL 文件
	sqlFile := "sqlStatement.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatal("无法读取 SQL 文件:", err)
	}
	sqlStatement := string(sqlBytes)

	wg := sync.WaitGroup{} // 创建WaitGroup用于等待所有goroutine执行完成
	// 打印日志
	log.Printf("更新数据")

	for _, keyword := range keywords {
		wg.Add(1) // 每个 goroutine 增加计数器
		go func(kw string, wg *sync.WaitGroup) {
			defer wg.Done() // 每个 goroutine 完成时减少计数器
			updateDatabase(db, sqlStatement, kw, wg)
		}(keyword, &wg)
	}
	log.Printf("等待所有更新操作完成")

	// 等待所有更新操作完成
	wg.Wait()
}
