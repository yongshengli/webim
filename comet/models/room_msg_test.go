package models

import (
    "testing"
    "github.com/astaxie/beego/orm"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // import your used driver
    "github.com/jinzhu/gorm"
    "time"
)

func TestRoomMsgInsert(t *testing.T){

    connectMysql()

    o := orm.NewOrm()
    roomMsg := RoomMsg{RoomId:"1", Content:"ssssss", Uid:123}
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
    db, _ := connectMysql()
    arr, _ := FindRoomMsgLast(db, "sss", 3)
    fmt.Println(arr)

    data := &RoomMsg{RoomId:"sss", Uid:123, Content:"dsdfsfdfsdfds", CT: time.Now().Unix()}
    res := db.Table("room_msg").Create(data)
    //res := InsertRoomMsg(db, "sss", data)
    fmt.Println(res.Error)
    //for _, row := range arr{
    //    fmt.Println(row)
    //}
}