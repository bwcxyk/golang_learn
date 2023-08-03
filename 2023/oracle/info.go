package main

import (
	"database/sql"
	"fmt"
	"sync"

	go_ora "github.com/sijms/go-ora/v2"
)

func Query(db *sql.DB, query string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 执行查询
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}(rows)

	// 处理查询结果
	for rows.Next() {
		// 读取查询结果
		var result string
		err := rows.Scan(&result)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		fmt.Println("Query result:", result)
	}
}

func main() {
	dsn := go_ora.BuildUrl("192.168.1.70", 1521, "orcl", "tms_user", "123456", nil)

	// 连接数据库
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing database:", err)
		}
	}(db)

	var wg sync.WaitGroup
	wg.Add(1)
	Query(db, "select sysdate from dual", &wg)
	wg.Wait()
}
