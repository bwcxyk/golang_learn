/*
@Author : yaokun
@Time : 2020/8/3 9:41
*/

package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)


func main() {
	conn, err := redis.Dial("tcp",
		"192.168.1.90:6379",
		redis.DialUsername(""),
		redis.DialPassword("redis"),
		redis.DialDatabase(1))
	if err != nil {
		fmt.Println("Connect to redis failed ,cause by >>>", err)
		return
	}
	defer conn.Close()

	// 写数据
	_, err = conn.Do("Set", "name", "zhangsan")
	if err != nil {
		fmt.Println("redis set value failed >>>", err)
		return
	}

	// 读数据
	data, err := redis.String(conn.Do("Get", "name"))
	if err != nil {
		fmt.Println("redis get value failed >>>", err)
		return
	}
	fmt.Println(data)

	// 删除数据
	_, err = conn.Do("DEL", "name")
	if err != nil {
		fmt.Println("redis.write err=", err)
		return
	}

}