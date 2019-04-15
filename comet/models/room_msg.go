package models

import (
    "comet/src/rrx/farm"
    "fmt"
    "github.com/astaxie/beego/orm"
)
const ROOM_MSG_NUM = 1024

type RoomMsg struct {
    Id      int    `orm:"auto"`
    RoomId  int
    Uid     int
    Content string `orm:"size(100)"`
    CT      int    `orm:"auto_now_add;type(datetime)"`
}


func RoomMsgTableName(roomId string) string {
    h := farm.Hash32([]byte(roomId))
    return fmt.Sprintf("room_msg_%d", int(h)%ROOM_MSG_NUM)
}

func FindRoomMsgLast(roomId string, limit int) (total int64, arr []RoomMsg, err error) {
    arr = []RoomMsg{}
    o := orm.NewOrm()
    sqlStr := fmt.Sprintf(`select * from %s where room_id=? order by id desc limit ?`, RoomMsgTableName(roomId))
    total, err = o.Raw(sqlStr,  roomId, limit).QueryRows(&arr)
    if err!=nil{
        return
    }
    return
}
func InsertRoomMsg(roomId, data map[string]interface{}){
    o := orm.NewOrm()
    o.Raw(``).Exec()
}