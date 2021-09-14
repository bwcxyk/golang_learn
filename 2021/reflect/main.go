/*
@Author : yaokun
@Time : 2020/8/14 9:57
*/

package main

import (
	"fmt"
	"reflect"
)

func main() {
	a :="1"
	fmt.Println("a:",reflect.TypeOf(a))
	a1 := a[0]
	fmt.Println("a1:",reflect.TypeOf(a1))
	fmt.Println("a1:" ,a1)
	a2 :='0'
	fmt.Println("a2:",reflect.TypeOf(a2))
	fmt.Println("a2:" ,a2)
	a3 := a[0] - '0'
	fmt.Println("a3:",reflect.TypeOf(a3))
	fmt.Println("a3:" ,a3)
}