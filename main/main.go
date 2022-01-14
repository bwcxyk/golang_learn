/*
@Author : yaokun
@Time : 2020/7/16 11:06
*/
package main

import (
	"fmt"
	"main/calc"
	"main/dance"
)

func main() {
	a :=1
	b :=1
	a =2
	res :=calc.Add(a,b)
	fmt.Printf("%d + %d = %d \n", a ,b ,res)
	fmt.Println("hello world")
	dance.WhoDance()
}
