package models

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type RoomUser struct {
	Id     uint64 `json:"id" gorm:"primary_key"`
	RoomId uint64 `json:"room_id"`
	Uid    uint64 `json:"uid"`
	CT     int64  `json:"c_t" gorm:"column:c_t"`
}

/**
 * 向room 中插入用户
 */
func InserRoomUser(ru *RoomUser) *gorm.DB {
	if db.NewRecord(ru) == false {
		return nil
	}
	if ru.CT < 1 {
		ru.CT = time.Now().Unix()
	}
	return db.Table("room_user").Create(ru)
}
func FindByRoomId(roomId uint64) ([]RoomUser, *gorm.DB) {
	var list []RoomUser
	res := db.Table("room_user").Where("room_id=?", roomId).Find(&list)
	return list, res
}

/*
 * 批量添加用户到聊天室
 */
func BatchInserRoomUser(roomId uint64, uids []uint64) *gorm.DB {
	slice := make([]interface{}, 0)
	ct := time.Now().Unix()
	tpl := "insert into `room_user` (`room_id`,`uid`,`c_t`) values "
	for _, uid := range uids {
		tpl += `(?),`
		slice = append(slice, []interface{}{roomId, uid, ct})
	}
	tpl = strings.Trim(tpl, ",")
	//db.LogMode(true)
	return db.Exec(tpl, slice...)
}
