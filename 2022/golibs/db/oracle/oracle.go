/*
@Author : YaoKun
@Time : 2022/9/30 10:15
*/

package oracle

import (
	"database/sql"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/sijms/go-ora/v2"
	"golibs/config"
	"golibs/global/consts"
	"log"
	"sync"
)

type GoLibOracle = *xorm.EngineGroup

var oracleOnce sync.Once
var oracleNew GoLibOracle

func NewOracle() GoLibOracle {
	oracleOnce.Do(func() {
		var err error
		dsn := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
			config.Config.GetOracle().GetUsername(),
			config.Config.GetOracle().GetPassword(),
			config.Config.GetOracle().GetHost(),
			config.Config.GetOracle().GetPort(),
			config.Config.GetOracle().GetDbname(),
		)
		oracleNew, err = xorm.NewEngineGroup("oracle", []string{dsn})
		if config.Config.GetActive() == consts.EnvDev {
			oracleNew.ShowSQL(true)
		}
		if err != nil {
			panic(err)
		}
	})
	return oracleNew
}

func TestConnect() {
	osqlInfo := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		config.Config.GetOracle().GetUsername(),
		config.Config.GetOracle().GetPassword(),
		config.Config.GetOracle().GetHost(),
		config.Config.GetOracle().GetPort(),
		config.Config.GetOracle().GetDbname(),
	)
	db, err := sql.Open("oracle", osqlInfo)
	if err != nil {
		log.Fatalf("connect oracle db error: %s:", err.Error())
	}

	rows, err := db.Query("select to_char(sysdate,'yyyy-mm-dd hh24:mi:ss') AS name from dual")
	if err != nil {
		fmt.Println("exec query error:", err.Error())
	}
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Println("exec query error:", err.Error())
		}
		fmt.Println("fetch item:")
		fmt.Println(name)
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("1234")
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println("链接成功")
}

