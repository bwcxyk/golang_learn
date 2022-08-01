/*
@Author : YaoKun
@Time : 2021/9/10 10:21
*/

package config

import (
	"fmt"
	"github.com/spf13/viper"
)

// 读取配置文件config

type Config struct {
	Redis     string
	MySQL     MySQLConfig
	HuaWeiOcr HuaWeiOcrConfig
}

type MySQLConfig struct {
	Port     int
	Host     string
	Username string
	Password string
}

type HuaWeiOcrConfig struct {
	AccessKey       string
	SecretAccessKey string
}

func InitConfig() {
	// 把配置文件读取到结构体上
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	viper.Unmarshal(&config) //将配置文件绑定到config上
	//fmt.Println("config: ", config.MySQL.Username)
	//fmt.Println("config: ", config.HuaWeiOcr.AccessKey)
	//fmt.Println("config: ", config.Redis)
}
