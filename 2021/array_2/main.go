/*
@Author : YaoKun
@Time : 2022/1/20 15:51
*/

package main

import "fmt"

func main() {
	var a = [...]int{1, 3, 5, 7, 8}
	var sum int
	for i := 0; i < len(a); i++ {
		fmt.Println("-------a[i]-------", a[i])
		sum += a[i]
		// len(a)-1值为4
		if i == len(a)-1 {
			fmt.Println("-------sum-------", sum)
		}
	}
}
