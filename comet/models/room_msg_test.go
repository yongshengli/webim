package models

import (
    "testing"
    "github.com/astaxie/beego/orm"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // import your used driver
)

func initMysql(){
    orm.RegisterModel(new(RoomMsg))
    // set default database
    orm.RegisterDataBase("default", "mysql", "root:123456@/chat?charset=utf8", 30)

    // create table
    orm.RunSyncdb("default", false, true)
}
func TestRoomMsgInsert(t *testing.T){

    initMysql()

    o := orm.NewOrm()
    roomMsg := RoomMsg{RoomId:1, Content:"ssssss", Uid:123}
    res, err := o.Insert(&roomMsg)
    if err != nil {
        t.Error(err)
    }
    fmt.Println(res)
}
func TestRoomMsgTableName(t *testing.T) {
    fmt.Println(RoomMsgTableName("0"))
}
func TestFindRoomMsgLast(t *testing.T) {
    initMysql()
    total, arr, err:= FindRoomMsgLast("sss", 12)
    if err!=nil{
        fmt.Println(err)
    }
    fmt.Println(total, arr)
}