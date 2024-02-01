package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/gen2brain/dlgs"
	"github.com/xuri/excelize/v2"
)

// 下载函数，从指定URL下载文件并保存为指定文件名
func download(fileName, url string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{} // 获取信号量

	// 发送HTTP GET请求到指定URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s 下载失败：%v\n", fileName, err)
		<-sem // 释放信号量
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s 下载失败，状态码：%d\n", fileName, resp.StatusCode)
		<-sem // 释放信号量
		return
	}

	// 从URL中提取文件扩展名
	fileExtension := path.Ext(url)
	fileNameWithExtension := fileName + fileExtension

	// 创建目标文件并保存下载的内容
	file, err := os.Create(fileNameWithExtension)
	if err != nil {
		fmt.Printf("%s 创建文件失败：%v\n", fileNameWithExtension, err)
		<-sem // 释放信号量
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("%s 保存文件失败：%v\n", fileNameWithExtension, err)
	}

	<-sem // 释放信号量
	fmt.Printf("%s 下载完成\n", fileNameWithExtension)
}

func processExcelFile(filename string) {
	// 打开Excel文件
	excelFile, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Printf("打开Excel文件失败：%v\n", err)
		return
	}

	// 获取第一个工作表
	sheetName := "Sheet1" // 指定工作表名称
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		fmt.Printf("读取工作表失败：%v\n", err)
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 4) // 创建信号量，限制并发数为 4

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]

		if rowIndex == 0 {
			continue // 跳过标题行
		}

		fileNameBase := row[0]
		for colIndex, columnValue := range row[1:] {
			fileName := fmt.Sprintf("%s_%d", fileNameBase, colIndex+1)
			url := columnValue

			wg.Add(1)
			go download(fileName, url, &wg, sem)
		}
	}

	wg.Wait()
}

func main() {
	filePath, _, err := dlgs.File("选择Excel文件", "*.xlsx", false)
	if err != nil {
		fmt.Println("文件选择错误:", err)
		return
	}

	if filePath == "" {
		fmt.Println("没有选择文件")
		return
	}

	processExcelFile(filePath)
}
