/*
   Copyright (c) [2021] itsos
   golibs is licensed under Mulan PSL v2.
   You can use this software according to the terms and conditions of the Mulan PSL v2.
   You may obtain a copy of Mulan PSL v2 at:
               http://license.coscl.org.cn/MulanPSL2
   THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
   See the Mulan PSL v2 for more details.
*/

package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"golibs/global/variable"
	"golibs/utils/reflects"
	"reflect"
	"strconv"
	"strings"
)

// ConfigurationReadOnly web 应用基础配置
type ConfigurationReadOnly interface {
	GetActive() string
	GetUrl() string
	GetSwaggerUrl() string
	GetDomain() string
	GetPort() string
	GetTimezone() string
	GetSock() string
	GetSwaggerPort() string
	GetScheme() string
	GetLogFile() string
	GetIsSetupView() bool
	GetRedisUse() string
	GetCryptAesToken() string
	GetCryptRsaPriv() string
	GetCrosAllowOrigin() []string
	GetSignatureExclude() []string
	GetCrosAllowHeaders() string
	GetEs() []string
	GetMysql() map[string]IMysql
	GetOracle() IOracle
	GetRedis() IRedis
}

type Configuration struct {
	Active           string `yaml:"active"`
	Domain           string `yaml:"domain"`
	Port             string `yaml:"port"`
	Scheme           string `yaml:"scheme"`
	Sock             string `yaml:"sock"`
	Timezone         string `yaml:"timezone"`
	SwaggerPort      string `yaml:"swagger.port"`
	Es               string `yaml:"es"`
	Mysql            string `yaml:"mysql"`
	Oracle           string `yaml:"oracle"`
	Redis            string `yaml:"redis"`
	RedisUse         string `yaml:"redis_use"`
	CryptAesToken    string `yaml:"crypt.aes.token"`
	CryptRsaPriv     string `yaml:"crypt.rsa.priv"`
	LogFile          string `yaml:"logfile"`
	IsSetupView      string `yaml:"is_setup_view"`
	CrosAllowOrigin  string `yaml:"cros.allow_origin"`
	CrosAllowHeaders string `yaml:"cros.allow_headers"`
	SignatureExclude string `yaml:"signature.exclude"`
}

func (c Configuration) GetUrl() string {
	url := c.GetScheme() + "://" + c.GetDomain()
	if c.GetPort() != "" {
		url += ":" + c.GetPort()
	}
	return url
}

func (c Configuration) GetSwaggerUrl() string {
	url := c.GetScheme() + "://" + c.GetDomain()
	if c.GetSwaggerPort() != "" {
		url += ":" + c.GetSwaggerPort()
	}
	return url
}

// GetEs 获取es配置
func (c Configuration) GetEs() []string {
	return viper.GetStringSlice(c.Es)
}

// GetSock 获取socket文件地址
func (c Configuration) GetSock() string {
	return viper.GetString(c.Sock)
}

func (c Configuration) GetTimezone() string {
	return viper.GetString(c.Timezone)
}

func (c Configuration) GetActive() string {
	return viper.GetString(c.Active)
}

func (c Configuration) GetRedisUse() string {
	return viper.GetString(c.RedisUse)
}

func (c Configuration) GetDomain() string {
	return viper.GetString(c.Domain)
}

func (c Configuration) GetPort() string {
	return viper.GetString(c.Port)
}

func (c Configuration) GetSwaggerPort() string {
	return viper.GetString(c.SwaggerPort)
}

func (c Configuration) GetScheme() string {
	return viper.GetString(c.Scheme)
}

func (c Configuration) GetLogFile() string {
	return viper.GetString(c.LogFile)
}

func (c Configuration) GetIsSetupView() bool {
	return viper.GetBool(c.IsSetupView)
}

func (c Configuration) GetCryptAesToken() string {
	return viper.GetString(c.CryptAesToken)
}

func (c Configuration) GetCryptRsaPriv() string {
	return viper.GetString(c.CryptRsaPriv)
}

func (c Configuration) GetCrosAllowOrigin() []string {
	return viper.GetStringSlice(c.CrosAllowOrigin)
}

func (c Configuration) GetCrosAllowHeaders() string {
	return strings.Join(viper.GetStringSlice(c.CrosAllowHeaders), ",")
}

func (c Configuration) GetSignatureExclude() []string {
	return viper.GetStringSlice(c.SignatureExclude)
}

// mysqlAlone mysql配置实例
var mysqlAlone map[string]IMysql

func (c Configuration) GetMysql() map[string]IMysql {
	if mysqlAlone == nil {
		mysqlAlone = make(map[string]IMysql)
		for k, v := range viper.GetStringMap(c.Mysql) {
			m := mysql{}
			if err := mapstructure.Decode(v, &m); err != nil {
				panic(err)
			}
			mysqlAlone[k] = m
		}
	}
	return mysqlAlone
}

// redisAlone redis配置实例
var redisAlone IRedis

func (c Configuration) GetRedis() IRedis {
	if redisAlone == nil {
		v := viper.GetStringMapString(c.Redis)
		r := redis{}
		if err := mapstructure.Decode(v, &r); err != nil {
			panic(err)
		}
		redisAlone = r
	}
	return redisAlone
}

// oracle配置实例
var oracleAlone IOracle

func (c Configuration) GetOracle() IOracle {
	if oracleAlone == nil {
		v := viper.GetStringMapString(c.Oracle)
		o := oracle{}
		if err := mapstructure.Decode(v, &o); err != nil {
			panic(err)
		}
		oracleAlone = o
	}
	return oracleAlone
}

// oracle配置
type oracle struct {
	Host     string
	//Port     int
	Port     string
	Username string
	Password string
	Dbname 	 string
}

func (o oracle) GetHost() string {
	return o.Host
}

//func (o oracle) GetPort() int {
//	return o.Port
//}

func (o oracle) GetPort() int {
	port, _ := strconv.Atoi(o.Port)
	return port
}

func (o oracle) GetUsername() string {
	return o.Username
}

func (o oracle) GetPassword() string {
	return o.Password
}

func (o oracle) GetDbname() string {
	return o.Dbname
}


type IOracle interface {
	GetHost()     string
	GetPort() 	  int
	GetUsername() string
	GetPassword() string
	GetDbname()   string
}

var _ IOracle = (*oracle)(nil)

// mysql mysql配置
type mysql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Charset  string
}

func (m mysql) GetHost() string {
	return m.Host
}

func (m mysql) GetPort() int {
	return m.Port
}

func (m mysql) GetUser() string {
	return m.User
}

func (m mysql) GetPassword() string {
	return m.Password
}

func (m mysql) GetDatabase() string {
	return m.Database
}

func (m mysql) GetCharset() string {
	return m.Charset
}

type IMysql interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetDatabase() string
	GetCharset() string
}

var _ IMysql = (*mysql)(nil)

// sqlite sqlite配置
type sqlite struct {
	Timezone    string `yaml:"sqlite.timezone"`
	StorageFile string `yaml:"sqlite.storage_file"`
}

func (s sqlite) GetTimezone() string {
	return viper.GetString(s.Timezone)
}

func (s sqlite) GetStorageFile() string {
	return viper.GetString(s.StorageFile)
}

type ISqlite interface {
	GetTimezone() string
	GetStorageFile() string
}

var _ ISqlite = (*sqlite)(nil)

// sqliteIns sqlite实例
var sqliteIns ISqlite = nil

func GetSqlite() ISqlite {
	if sqliteIns == nil {
		sqliteIns = &sqlite{}
		reflects.TagToValueFlip(reflect.ValueOf(sqliteIns).Elem(), reflects.YAML)
	}
	return sqliteIns
}

// redis 单机redis配置
type redis struct {
	Host     string
	Port     string
	Username string
	Password string
	Db       string
}

func (r redis) GetHost() string {
	return r.Host
}

func (r redis) GetPort() int {
	port, _ := strconv.Atoi(r.Port)
	return port
}

func (r redis) GetUsername() string {
	return r.Username
}

func (r redis) GetPassword() string {
	return r.Password
}

func (r redis) GetDb() int {
	db, _ := strconv.Atoi(r.Db)
	return db
}

type IRedis interface {
	GetHost() string
	GetPort() int
	GetUsername() string
	GetPassword() string
	GetDb() int
}

var _ IRedis = (*redis)(nil)

var _ ConfigurationReadOnly = (*Configuration)(nil)

// CovertConfiguration struct tag 指定为属性值，方便操作
func CovertConfiguration() *Configuration {
	loadConfigFile()
	c := &Configuration{}
	t := reflect.ValueOf(c).Elem()
	reflects.TagToValueFlip(t, reflects.YAML)
	return c
}

// loadConfigFile 加载配置文件
func loadConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(variable.ConfPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

var Config = CovertConfiguration()

// minio
type minio struct {
	Endpoint        string `yaml:"minio.endpoint"`
	AccessKeyID     string `yaml:"minio.access_key_id"`
	SecretAccessKey string `yaml:"minio.secret_access_key"`
	UseSSL          string `yaml:"minio.use_ssl"`
}

func (m minio) GetEndpoint() string {
	return viper.GetString(m.Endpoint)
}

func (m minio) GetAccessKeyID() string {
	return viper.GetString(m.AccessKeyID)
}

func (m minio) GetSecretAccessKey() string {
	return viper.GetString(m.SecretAccessKey)
}

func (m minio) GetUseSSL() bool {
	return viper.GetBool(m.UseSSL)
}

type IMinio interface {
	GetEndpoint() string
	GetAccessKeyID() string
	GetSecretAccessKey() string
	GetUseSSL() bool
}

var minioIns IMinio = nil

func GetMinio() IMinio {
	if minioIns == nil {
		minioIns = &minio{}
		reflects.TagToValueFlip(reflect.ValueOf(minioIns).Elem(), reflects.YAML)
	}
	return minioIns
}

// redisCluster
type redisCluster struct {
	Hosts    string `yaml:"redis_cluster.hosts"`
	Username string `yaml:"redis_cluster.username"`
	Password string `yaml:"redis_cluster.password"`
}

func (r redisCluster) GetHosts() []string {
	return viper.GetStringSlice(r.Hosts)
}

func (r redisCluster) GetUsername() string {
	return viper.GetString(r.Username)
}

func (r redisCluster) GetPassword() string {
	return viper.GetString(r.Password)
}

type IRedisCluster interface {
	GetHosts() []string
	GetUsername() string
	GetPassword() string
}

var redisClusterIns IRedisCluster = nil

func GetRedisCluster() IRedisCluster {
	if redisClusterIns == nil {
		redisClusterIns = &redisCluster{}
		reflects.TagToValueFlip(reflect.ValueOf(redisClusterIns).Elem(), reflects.YAML)
	}
	return redisClusterIns
}
