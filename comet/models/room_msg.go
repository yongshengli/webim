package models

import (
    "github.com/dgryski/go-farm"
    "fmt"
    "github.com/jinzhu/gorm"
)
const ROOM_MSG_NUM = 1024

type RoomMsg struct {
    Id      uint `gorm:"primary_key"`
    RoomId  string
    Uid     int
    Content string
    CT      int `gorm:"column:c_t"`
}

func RoomMsgTableName(roomId string) string {
    h := farm.Hash32([]byte(roomId))
    return fmt.Sprintf("room_msg_%d", int(h)%ROOM_MSG_NUM)
}
func (rm *RoomMsg) GetTableName(roomId string) string {
    return RoomMsgTableName(rm.RoomId)
}

func FindRoomMsgLast(db *gorm.DB, roomId string, limit int) (total int64, arr []RoomMsg, err error) {
    db.Table(RoomMsgTableName(roomId)).
        Where("room_id=?", roomId).
        Order("id desc").
        Limit(limit)
    if err!=nil{
        return
    }
    return
}

func InsertRoomMsg(db *gorm.DB, roomId string, data *RoomMsg) *gorm.DB{
    return db.Table(RoomMsgTableName(roomId)).Create(data)

}