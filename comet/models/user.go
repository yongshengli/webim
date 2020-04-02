package models

import (
	"time"
	"webim/comet/common"

	"github.com/jinzhu/gorm"
)

//User 用户模型
type User struct {
	Id        uint64 `json:"id" gorm:"primary_key"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CT        int64  `json:"c_t" gorm:"column:c_t"`
	LastLogin uint64 `json:"last_login"`
	DeviceId  string `json:"device_id"`
}

//FindByName 根据用户名查询用户记录
func FindByName(name string) (*User, error) {
	user := &User{}
	res := db.Table("user").Where("username=?", name).First(user)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}

//CheckPwd 检查密码是否正确
func CheckPwd(u *User, passwd string) bool {
	// return strings.EqualFold(u.Password, common.Md5(passwd))
	return u.Password == common.Md5(passwd)
}

//InsertUser 注册用户
func InsertUser(u *User) *gorm.DB {
	if db.NewRecord(u) == false {
		return nil
	}
	u.CT = time.Now().Unix()
	return db.Table("user").Create(u)
}
