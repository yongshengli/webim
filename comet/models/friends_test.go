package models

import (
	"testing"
)

func TestAddFriend(t *testing.T) {
	res := AddFriend(1, 2)
	if res.Error != nil {
		t.Error(res.Error)
	}
	res = AddFriend(1, 3)
	if res.Error != nil {
		t.Error(res.Error)
	}
	list, err := FindFriends(1)
	if err != nil {
		t.Error(err)
	}
	if len(list) < 2 {
		t.Error("添加或者查询用户好友失败")
	}
}
