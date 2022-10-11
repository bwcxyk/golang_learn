package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Get 发起get请求
func Get(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicf("An error occurred while requesting: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicf("An error occurred while closing the body: %v", err)
		}
	}(resp.Body)


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("An error occurred while reading data: %v", err)
	}
	return body
}

// PostJson 发起post请求
// params, _ := json.Marshal(map[string]string{"name": "Test"})
func PostJson(url string, params []byte) []byte {
	responseBody := bytes.NewBuffer(params)
	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Panicf("An error occurred while requesting: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicf("An error occurred while closing the body: %v", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("An error occurred while reading data: %v", err)
	}
	return body
}
