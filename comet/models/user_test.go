package models

import (
	"testing"
	"webim/comet/common"
)

func TestInsertUser(t *testing.T) {
	user := &User{UserName: "admin", Password: common.Md5("123456")}
	res := InsertUser(user)
	if res.Error != nil {
		t.Error(res)
	}
	if user.Id < 1 {
		t.Errorf("插入user失败 get uid:%d\n", user.Id)
	}
	if CheckPwd(user, "123456") != true {
		t.Errorf("验证密码接口错误 except %s but got %s\n", "true", "false")
	}
}
