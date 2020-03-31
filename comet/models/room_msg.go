package models

import (
	"fmt"
	"time"

	"github.com/dgryski/go-farm"
	"github.com/jinzhu/gorm"
)

const ROOM_MSG_NUM = 1024

type RoomMsg struct {
	Id      uint64 `json:"id" gorm:"primary_key"`
	RoomId  string `json:"room_id"`
	Uid     int64  `json:"uid"`
	Uname   string `json:"uname"`
	Content string `json:"content"`
	CT      int64  `json:"c_t" gorm:"column:c_t"`
}

/*
* RoomMsgTableName
* 获取room_msg 表名
 */
func RoomMsgTableName(roomId string) string {
	//暂时先使用一个表
	return "room_msg"
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
	if err != nil {
		return
	}
	return
}

//写入一条聊天数据数据
func InsertRoomMsg(roomId string, data *RoomMsg) *gorm.DB {
	if db.NewRecord(data) == false {
		return nil
	}
	if data.CT < 1 {
		data.CT = time.Now().Unix()
	}
	return db.Table(RoomMsgTableName(roomId)).Create(data)
}
