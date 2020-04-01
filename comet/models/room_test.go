package models

import (
	"testing"
)

var uids = []uint64{
	1,
	2,
	3,
	4,
	5,
}

func TestCreateRoom(t *testing.T) {
	// testData := &Room{RoomName: "dddd"}
	roomId := CreateRoom(uids)

	if roomId < 1 {
		t.Error("创建房间失败")
	}
	list, res := FindByRoomId(roomId)
	if res.Error != nil {
		t.Error(res.Error)
	}
	if len(list) < len(uids) {
		t.Error("写入数据条数不一致")
	}
	db.Table("room").Delete(&Room{}, "room_id=?", roomId)
}
func init() {
	ConnectTestMysql()
}
