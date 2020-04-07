package models

import (
	"testing"
)

func TestAddFriend(t *testing.T) {
	uid := uint64(1)
	res := AddFriend(uid, 2)
	if res.Error != nil {
		t.Error(res.Error)
	}
	res = AddFriend(uid, 3)
	if res.Error != nil {
		t.Error(res.Error)
	}
	list, err := FindFriends(uid)
	if err != nil {
		t.Error(err)
	}
	if len(list) < 2 {
		t.Error("添加或者查询用户好友失败")
	}
	db.Delete(&Friend{}, "uid=?", uid)
}
