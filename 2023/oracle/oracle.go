/*
@Author : YaoKun
@Time : 2023/8/3 13:42
*/
package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	go_ora "github.com/sijms/go-ora/v2"
)

func Query(ctx context.Context, db *sql.DB, query string, wg *sync.WaitGroup) error {
	defer wg.Done()

	// 使用context执行查询
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error closing rows:", err)
		}
	}(rows)

	// 处理查询结果
	for rows.Next() {
		var result string
		err := rows.Scan(&result)
		if err != nil {
			fmt.Println("读取查询结果时出错:", err)
			continue
		}
		fmt.Println("查询结果:", result)
	}

	return nil
}

func main() {
	dsn := go_ora.BuildUrl("192.168.1.70", 1521, "orcl", "tms_user", "123456", nil)

	// 连接数据库
	db, err := sql.Open("oracle", dsn)
	if err != nil {
		fmt.Println("连接数据库时出错:", err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing db:", err)
		}
	}(db) // 当主函数执行完毕时关闭数据库连接

	var wg sync.WaitGroup
	wg.Add(1)

	// 使用context设置查询超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 当主函数执行完毕时确保context被取消

	// 使用context运行查询
	err = Query(ctx, db, "select sysdate from dual", &wg)
	if err != nil {
		fmt.Println("执行查询时出错:", err)
	}

	wg.Wait()
}
