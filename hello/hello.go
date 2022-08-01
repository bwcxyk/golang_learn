/*
@Author : YaoKun
@Time : 2020/7/15 18:09
*/

package main

import "fmt"

func Hello() string {
	return "Hello, world"
}

func abc() int {
	return 123
}

func main() {
	fmt.Println(Hello())
	fmt.Println(abc())
}
