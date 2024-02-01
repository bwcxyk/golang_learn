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
    excelFile, err := excelize.OpenFile(filename)
    if err != nil {
        fmt.Printf("打开Excel文件失败：%v\n", err)
        dlgs.Error("错误", "打开Excel文件失败："+err.Error())
        return
    }

    sheetName := "Sheet1" // 指定工作表名称
    sheets := excelFile.GetSheetList()

    // 检查Sheet1是否存在
    sheetExists := false
    for _, name := range sheets {
        if name == sheetName {
            sheetExists = true
            break
        }
    }

    if !sheetExists {
        fmt.Printf("未找到工作表：%s\n", sheetName)
        dlgs.Error("错误", fmt.Sprintf("未找到工作表：%s", sheetName))
        return
    }

    rows, err := excelFile.GetRows(sheetName)
    if err != nil {
        fmt.Printf("读取工作表%s失败：%v\n", sheetName, err)
        dlgs.Error("错误", "读取工作表失败："+err.Error())
        return
    }
    // 检查是否有数据
    if len(rows) == 0 {
        fmt.Println("工作表是空的，没有数据可读取。")
        dlgs.Error("提示", "工作表是空的，没有数据可读取。")
        // 这里可以添加任何适当的逻辑来处理空工作表的情况
        return
    }

    var wg sync.WaitGroup
    sem := make(chan struct{}, 4) // 创建信号量，限制并发数为 4

    for rowIndex, row := range rows {
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
        dlgs.Error("文件选择错误", err.Error())
        return
    }

    if filePath == "" {
        fmt.Println("没有选择文件")
        return
    }

    processExcelFile(filePath)
}
