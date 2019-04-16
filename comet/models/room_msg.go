package models

import (
    "github.com/dgryski/go-farm"
    "fmt"
    "github.com/jinzhu/gorm"
    "time"
)
const ROOM_MSG_NUM = 1024

type RoomMsg struct {
    Id      uint `gorm:"primary_key"`
    RoomId  string
    Uid     int
    Content string
    CT      int64 `gorm:"column:c_t"`
}

func RoomMsgTableName(roomId string) string {
    h := farm.Hash32([]byte(roomId))
    return fmt.Sprintf("room_msg_%d", int(h)%ROOM_MSG_NUM)
}
func (rm *RoomMsg) GetTableName(roomId string) string {
    return RoomMsgTableName(rm.RoomId)
}
//获取最新的几条聊天室聊天记录
func FindRoomMsgLast(roomId string, limit int) (arr []RoomMsg, err error) {
    //arr = []RoomMsg{}
    res := db.Table(RoomMsgTableName(roomId)).
        Where("room_id=?", roomId).
        Order("id desc").
        Limit(limit).Find(&arr)
    err = res.Error
    if err!=nil{
        return
    }
    return
}

//写入一条聊天数据数据
func InsertRoomMsg(roomId string, data *RoomMsg) *gorm.DB{
    if db.NewRecord(data) == false {
        return nil
    }
    if data.CT < 1 {
        data.CT = time.Now().Unix()
    }
    return db.Table(RoomMsgTableName(roomId)).Create(data)
}