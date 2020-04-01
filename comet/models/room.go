package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Room room
type Room struct {
	RoomId   uint64 `json:"id" gorm:"primary_key"`
	RoomName string `json:"room_name"`
	CT       int64  `json:"c_t" gorm:"column:c_t"`
}

//InsertRoom 房间入表
func InsertRoom(r *Room) *gorm.DB {
	if db.NewRecord(r) == false {
		return nil
	}
	if r.CT < 1 {
		r.CT = time.Now().Unix()
	}
	return db.Table("room").Create(r)
}

//CreateRoom 创建房间并添加用户
func CreateRoom(uids []uint64) uint64 {
	ct := time.Now().Unix()
	transaction := db.Begin()
	rData := &Room{
		RoomName: "room",
		CT:       ct,
	}
	res := InsertRoom(rData)
	if res.Error != nil {
		transaction.Callback()
		return 0
	}
	res = BatchInserRoomUser(rData.RoomId, uids)
	if res.Error != nil {
		transaction.Callback()
		return 0
	}
	transaction.Commit()

	return rData.RoomId
}
