package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	downloadDestFolder = "E:/file"
	urlFilePath        = "E:/file/file.txt"
)

func init() {
	log.SetFlags(log.Lshortfile)
	_ = os.MkdirAll(downloadDestFolder, 0777)
}

func main() {
	fi, err := os.Open(urlFilePath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	var w sync.WaitGroup
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			log.Println("read url complete")
			break
		}
		list := strings.Split(string(line), ",")
		w.Add(1)
		go download(list[1], list[0]+".xlsx", &w)
	}
	w.Wait()

}

func download(url string, filename string, w *sync.WaitGroup) {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get -> %v", err)
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll -> %s", err.Error())
		return
	}
	defer res.Body.Close()
	if err = ioutil.WriteFile(downloadDestFolder+string(filepath.Separator)+filename, data, 0777); err != nil {
		log.Println("Error Saving:", filename, err)
	} else {
		log.Println("Saved:", filename)
	}
	w.Done()
}
