package models

import (
	"comet/common"
	"fmt"
	"time"

	"github.com/dgryski/go-farm"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
)

const ROOM_MSG_NUM = 1024

type RoomMsg struct {
	Id      uint64 `json:"id" gorm:"primary_key"`
	MsgId   uint64 `json:"msg_id"`
	RoomId  string `json:"room_id"`
	Uid     int64  `json:"uid"`
	Uname   string `json:"uname"`
	Content string `json:"content"`
	CT      int64  `json:"c_t" gorm:"column:c_t"`
}

//RoomMsgTableName 获取room_msg 表名
func RoomMsgTableName(roomId string) string {
	//暂时先使用一个表
	return "room_msg"
	h := farm.Hash32([]byte(roomId))
	return fmt.Sprintf("room_msg_%d", int(h)%ROOM_MSG_NUM)
}

//GetTableName
func (rm *RoomMsg) GetTableName(roomId string) string {
	return RoomMsgTableName(rm.RoomId)
}

//FindRoomMsgLast 获取最新的几条聊天室聊天记录
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

//GetLastRoomMsg 获取最近的消息
func GetLastRoomMsg(roomId string, limit int) ([]RoomMsg, error) {

	byteArr, err := redis.ByteSlices(common.RedisClient.Do("zrevrange", msgZsetKey(roomId), 0, limit))
	if err != nil {
		return nil, err
	}
	list := make([]RoomMsg, 0, len(byteArr))
	for _, r := range byteArr {
		var tmp RoomMsg
		common.DeJson(r, &tmp)
		list = append(list, tmp)
	}
	return list, nil
}

//CreateMsgId 生成消息id
func CreateMsgId(roomId string) (uint64, error) {
	return common.RedisClient.Incr(msgIdKey(roomId))
}

//SaveMsg2Redis 保存消息到redis中
func SaveMsg2Redis(data *RoomMsg) (int64, error) {
	var err error
	roomId := data.RoomId
	data.MsgId, err = CreateMsgId(roomId)
	if err != nil {
		return 0, err
	}
	if data.CT < 1 {
		data.CT = time.Now().Unix()
	}
	d, _ := common.EnJson(data)
	reply, err := redis.Int64(common.RedisClient.Do("zadd", msgZsetKey(roomId), data.MsgId, d))
	if err != nil {
		return reply, err
	}
	common.RedisClient.Do("ZREMRANGEBYRANK", msgZsetKey(roomId), 0, -3000)
	return common.RedisClient.Expire(msgZsetKey(roomId), 7*24*time.Hour)
}

func msgIdKey(roomId string) string {
	return fmt.Sprintf("room:%s:msgId", roomId)
}
func msgZsetKey(roomId string) string {
	return fmt.Sprintf("room:%s:msgZset", roomId)
}

//InsertRoomMsg 写入一条聊天数据数据
func InsertRoomMsg(data RoomMsg) *gorm.DB {
	roomId := data.RoomId
	if roomId == "" {
		return nil
	}
	// if db.NewRecord(data) == false {
	// 	return nil
	// }
	if data.CT < 1 {
		data.CT = time.Now().Unix()
	}
	return db.Table(RoomMsgTableName(roomId)).Create(&data)
}
