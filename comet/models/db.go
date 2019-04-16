package models

import (
    "github.com/jinzhu/gorm"
    "github.com/astaxie/beego/logs"
    "fmt"
)

var db *gorm.DB

func ConnectMysql(user, pass, dbName string) error {
    if db != nil {
        return nil
    }
    var err error
    db, err = gorm.Open("mysql", "mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8", user, pass, dbName))
    if err != nil {
        logs.Error("msg[链接mysql错误] err[%s]", err)
        return err
    }
    db.DB().SetMaxIdleConns(10)
    return nil
}

func ConnectTestMysql() (*gorm.DB, error){
    err := ConnectMysql("root", "123456", "chat")
    return db, err
}