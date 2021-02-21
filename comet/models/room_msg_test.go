package models

import (
	"fmt"
	"testing"
	"time"
)

func TestMsgIdKey(t *testing.T) {
	msgId := msgIdKey("3")
	if msgId != "room:3:msgId" {
		t.Error("msgIdKey err")
	}
}
func TestRoomMsgInsert(t *testing.T) {

	msgId, err := CreateMsgId("1")
	if err != nil {
		t.Error(err)
	}
	roomMsg := RoomMsg{MsgId: msgId, RoomId: "1", Content: "sss", Uid: 123, CT: time.Now().Unix()}
	res := InsertRoomMsg(roomMsg)
	if res.RowsAffected < 1 {
		t.Error("插入失败")
	}
	if db.NewRecord(roomMsg) {
		t.Error("插入返回的主键为空")
	}
	fmt.Println(roomMsg.Id)
}

func TestRoomMsgTableName(t *testing.T) {
	tableName := RoomMsgTableName("0")
	if tableName != "room_msg" {
		t.Errorf("table name err, except room_msg_802 but get %s", tableName)
	}
}

func TestFindRoomMsgLast(t *testing.T) {
	roomId := "100"
	arr, _ := FindRoomMsgLast(roomId, 3)
	num := len(arr)
	data := &RoomMsg{RoomId: roomId, Uid: 123, Content: "dsdfsfdfsdfds", CT: time.Now().Unix()}
	res := db.Table("room_msg").Create(data)
	if res.Error != nil {
		t.Error(res.Error)
	}
	arr2, _ := FindRoomMsgLast(roomId, 3)
	if len(arr2) <= num {
		t.Error("插入roomMsg数据失败")
	}
	db.Table("room_msg").Delete(&RoomMsg{}, "room_id=?", roomId)
}

func TestSaveMsg2Redis(t *testing.T) {
	res, err := SaveMsg2Redis(&RoomMsg{RoomId: "3", Uid: 1, Uname: "admin", Content: "123323232"})
	if err != nil {
		t.Error(err)
	}
	if res < 1 {
		t.Error("设置超时时间失败")
	}
	list, err := GetLastRoomMsg("3", 100)
	if err != nil {
		t.Error(err)
	}
	if len(list) < 1 {
		t.Error("获取最近消息数为0 except 1")
	}
	fmt.Println(list)
}
