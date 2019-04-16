package models

import (
    "testing"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // import your used driver
    "time"
)

func TestRoomMsgInsert(t *testing.T){

    ConnectTestMysql()
    roomMsg := RoomMsg{RoomId:"1", Content:"sss", Uid:123, CT:time.Now().Unix()}
    res := InsertRoomMsg(roomMsg.RoomId, &roomMsg)
    if res.RowsAffected < 1 {
        t.Error("插入失败")
    }
    if db.NewRecord(roomMsg) {
        t.Error("插入返回的主键为空")
    }
    fmt.Println(roomMsg.Id)
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