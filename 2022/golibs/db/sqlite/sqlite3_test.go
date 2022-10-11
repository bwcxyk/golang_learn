package sqlite

import (
	"log"
	"testing"
	"time"
)

type Role struct {
	Id         int64 `xorm:"not null autoincr pk id"` //指定主键并自增
	Name       string
	Info       string
	UpdateTime time.Time `xorm:"not null updated"` //修改后自动更新时间
	CreateTime time.Time `xorm:"not null created"` //创建时间
}

func initConfig() GoLibSqlite {
	return NewSqlite()
}

func TestConnect(t *testing.T) {
	initConfig()
}

func TestCreateTable(t *testing.T) {
	x := initConfig()

	role := new(Role)
	err := x.Sync2(role)
	if err != nil {
		log.Println(err)
		return
	}
}

// add
func TestAdd(t *testing.T) {
	x := initConfig()
	role := new(Role)
	role.Name = "普通"
	role.Info = "超级普通"
	a, e := x.Insert(role)
	t.Log(a, e)
	if a == 1 {
		b, err := x.ID(1).Update(role)
		t.Log(b, err)
	}
}

func TestTime(t *testing.T) {
	t.Log(time.Now().Format("2006-01-02"))
}
