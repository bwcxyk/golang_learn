/*
@Author : YaoKun
@Time : 2021/9/10 10:30
*/

package main

import (
	"IdCard"
	"config"
)

func main() {
	config.InitConfig()
	IdCard.IdCard()
}