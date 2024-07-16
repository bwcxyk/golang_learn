package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tedcy/fdfs_client"
)

var client *fdfs_client.Client

// 初始化 FastDFS 客户端
func initClient() error {
	var err error
	client, err = fdfs_client.NewClientWithConfig("./config/fastdfs.conf")
	if err != nil {
		fmt.Println("打开fast客户端失败", err.Error())
	}
	return err
}

// 清理 FastDFS 客户端
func cleanupClient() {
	if client != nil {
		client.Destory()
	}
}

// 上传文件到 fastDFS 系统
func UploadFile(fileName string) string {
	err := initClient()
	if err != nil {
		return ""
	}
	defer cleanupClient()

	fileId, err := client.UploadByFilename(fileName)
	if err != nil {
		fmt.Println("上传文件失败", err.Error())
		return ""
	}
	return fileId
}

// 下载文件
func DownLoadFile(fileId, tempFile string) {
	err := initClient()
	if err != nil {
		return
	}
	defer cleanupClient()

	err = client.DownloadToFile(fileId, tempFile, 0, 0)
	if err != nil {
		fmt.Println("下载文件失败", err.Error())
	}
}

// 删除文件
func DeleteFile(fileId string) {
	err := initClient()
	if err != nil {
		return
	}
	defer cleanupClient()

	err = client.DeleteFile(fileId)
	if err != nil {
		fmt.Println("删除文件失败", err.Error())
	}
}

func main() {
	// 检查参数是否足够
	if len(os.Args) < 2 {
		fmt.Println("请提供要上传的文件夹路径作为参数")
		return
	}
	// 要上传的文件夹路径
	dir := os.Args[1]

	// 遍历文件夹中的所有文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 跳过文件夹
		if info.IsDir() {
			return nil
		}
		// 上传文件
		fileId := UploadFile(path)
		if fileId != "" {
			fmt.Printf("文件 %s 上传成功，文件ID为: %s\n", path, fileId)
		} else {
			fmt.Printf("文件 %s 上传失败\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("遍历文件夹出错:", err)
	}
}
