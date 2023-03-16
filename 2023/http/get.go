/*
@Author : YaoKun
@Time : 2023/3/14 16:19
*/

package main

import (
	"encoding/json"
	"net/http"
)

func main() {

	http.HandleFunc("/hello", GetHanle)
	http.ListenAndServe("0.0.0.0:8080", nil)

}

func GetHanle(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "hello world")
	w.Header().Set("content-type", "text/json")
	msg := make(map[string]string)
	msg["code"] = "200"
	msg["msg"] = "success"
	ret, _ := json.Marshal(msg)
	w.Write(ret)
}
