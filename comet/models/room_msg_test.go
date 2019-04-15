package models

import (
    "testing"
    "github.com/astaxie/beego/orm"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // import your used driver
    "github.com/jinzhu/gorm"
)

func connectMysql() (*gorm.DB, error){
    // set default database
    db, err :=gorm.Open("mysql", "mysql", "root:123456@/chat?charset=utf8")
    if err!=nil{
        fmt.Println(err)
        return nil, err
    }
    db.DB().SetMaxIdleConns(10)
    return db, nil

}
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
    var arr []RoomMsg
    db.Table(RoomMsgTableName("sss")).
        Where("room_id=?", "sss").
        Order("id desc").
        Limit(3).Find(&arr)
    fmt.Println(arr)

    data := &RoomMsg{RoomId:"sss", Uid:123, Content:"dsdfsfdfsdfds"}
    res := InsertRoomMsg(db, "sss", data)
    fmt.Println(res.Error)
    //for _, row := range arr{
    //    fmt.Println(row)
    //}
}