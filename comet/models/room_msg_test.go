package models

import (
    "testing"
    "github.com/astaxie/beego/orm"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // import your used driver
    "time"
)

func TestRoomMsgInsert(t *testing.T){

    ConnectTestMysql()

    o := orm.NewOrm()
    roomMsg := RoomMsg{RoomId:"1", Content:"ssssss", Uid:123}
    res, err := o.Insert(&roomMsg)
    if err != nil {
        t.Error(err)
    }
    fmt.Println(res)
}

func TestRoomMsgTableName(t *testing.T) {
    tableName := RoomMsgTableName("0")
    if tableName != "room_msg_802" {
        t.Errorf("table name err, except room_msg_802 but get %s", tableName)
    }
}

func TestFindRoomMsgLast(t *testing.T) {
    err := ConnectTestMysql()
    if err!=nil {
        t.Error(err)
    }
    arr, _ := FindRoomMsgLast("sss", 3)
    fmt.Println(arr)

    data := &RoomMsg{RoomId:"sss", Uid:123, Content:"dsdfsfdfsdfds", CT: time.Now().Unix()}
    res := db.Table("room_msg").Create(data)
    //res := InsertRoomMsg(db, "sss", data)
    fmt.Println(res.Error)
}