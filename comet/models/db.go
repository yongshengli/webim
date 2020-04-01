package models

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql" // import your used driver
	"github.com/jinzhu/gorm"
)

var (
	db       *gorm.DB
	deviceDb *gorm.DB
)

//ConnectMysql 连接数据库
func ConnectMysql(host, port, user, pass, dbName string) error {
	if db != nil {
		return nil
	}
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", user, pass, host, port, dbName))
	if err != nil {
		logs.Error("msg[链接mysql错误] err[%s]", err)
		return err
	}
	db.DB().SetMaxIdleConns(10)
	return nil
}

//ConnectTestMysql 连接测试数据库用于单元测试
func ConnectTestMysql() error {
	err := ConnectMysql("127.0.0.1", "3306", "root", "123456", "chat")
	return err
}
