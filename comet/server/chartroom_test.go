package server

import (
	"comet/common"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestNewRoom(t *testing.T) {
	roomId := "123"
	room, err := NewRoom(roomId, "")
	if err != nil {
		t.Error(err)
	}
	if room == nil {
		t.Fail()
	}
	room, err = GetRoom(roomId)
	if err != nil {
		t.Error(err)
	}
	if room.Id == "" {
		t.Fail()
	}
	_, err = room.Join(&Session{DeviceToken: uuid.NewV4().String(), IP: "127.0.0.1", User: &User{Id: "1", Name: "张三"}})
	if err != nil {
		t.Error(err)
	}
	users, err := room.Users(0, -1)
	if err != nil {
		t.Error(err)
	}
	if len(users) < 1 {
		t.Error("用户进入房间失败")
	}
	roomArr, err := RoomList(0, -1)
	if err != nil {
		t.Error(err)
	}
	if len(roomArr) < 1 {
		t.Error("room zset 中没有数据")
	}
	roomNum, err := TotalRoom()
	if err != nil {
		t.Error(err)
	}
	if roomNum < 1 {
		t.Error("房间数错误")
	}
	res, err := DelRoom(roomId)
	if err != nil {
		t.Error(err)
	}
	if res < 1 {
		t.Error("删除聊天室失败")
	}
	room, err = GetRoom(roomId)
	if err != nil {
		t.Error(err)
	}
}

func init() {
	common.RedisInitTest()
}
