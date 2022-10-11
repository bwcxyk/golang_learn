package mysql

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

func initConfig() GoLibMysql {
	db := NewMysql()
	return db
}

func TestInitConfig(t *testing.T) {
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

// delete
func TestDelete(t *testing.T) {
	x := initConfig()
	r, err := x.ID(1).Delete(new(Role))
	t.Log(r, err)
}

// update
func TestUpdate(t *testing.T) {
	x := initConfig()
	role := new(Role)
	role.Name = "普通0"
	role.Info = "超级普通0"
	r, err := x.ID(1).Update(role)
	t.Log(r, err)
}

// select one
func TestSelectOne(t *testing.T) {
	x := initConfig()
	role := new(Role)
	role.Name = "普通0"
	role.Info = "超级普通0"
	r, err := x.Table("role").ID(1).Get(role)
	t.Log(r, err)
	t.Log(role.Name)
}

// select more
func TestSelectMore(t *testing.T) {
	x := initConfig()
	role := make([]Role, 0)
	err := x.Find(&role)
	t.Log(err, role)
}

// transaction
func TestTransaction(t *testing.T) {
	x := initConfig()
	session := x.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Fatal(err)
	}
	role := new(Role)
	role.Name = "普通"
	role.Info = "超级普通"
	_, e := session.Insert(role)
	if e != nil {
		log.Fatal(e)
	}
	session.Commit()
}
