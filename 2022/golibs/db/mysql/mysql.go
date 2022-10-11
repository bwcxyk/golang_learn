package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"golibs/config"
	"golibs/global/consts"
	"sync"
	"time"
)

// https://github.com/go-xorm/xorm
// https://gobook.io/read/gitea.com/xorm/manual-zh-CN/chapter-01/2.engine_group.html#

type GoLibMysql = *xorm.EngineGroup

func getDsn(mysql config.IMysql) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		mysql.GetUser(),
		mysql.GetPassword(),
		mysql.GetHost(),
		mysql.GetPort(),
		mysql.GetDatabase(),
		mysql.GetCharset(),
	)
}

var mysqlOnce sync.Once
var mysqlNew GoLibMysql

func NewMysql() GoLibMysql {
	mysqlOnce.Do(func() {
		var dataSourceNameSlice []string
		for _, mysql := range config.Config.GetMysql() {
			dataSourceNameSlice = append(dataSourceNameSlice, getDsn(mysql))
		}
		var err error
		mysqlNew, err = xorm.NewEngineGroup("mysql", dataSourceNameSlice)
		if err != nil {
			panic(err)
		}
		mysqlNew.TZLocation, _ = time.LoadLocation(config.Config.GetTimezone())
		mysqlNew.DatabaseTZ, _ = time.LoadLocation(config.Config.GetTimezone())
		if config.Config.GetActive() == consts.EnvDev {
			mysqlNew.ShowSQL(true)
		}
	})
	return mysqlNew
}
