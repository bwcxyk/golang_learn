/*
@Author : YaoKun
@Time : 2021/10/28 10:38
*/

package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func main() {
	for i := 1; i < 10; i++ {
		for j := 1; j < i+1; j++ {
			fmt.Print(func() string {
				var buf bytes.Buffer
				err := template.Must(template.New("f").Parse("{{.i}}x{{.j}}={{.expr1}}")).
					Execute(&buf, map[string]interface{}{"i": i, "j": j, "expr1": i * j})
				if err != nil {
					panic(err)
				}
				return buf.String()
			}(), " ")
		}
		fmt.Println("")
	}
}
