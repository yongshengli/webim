package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Friend struct {
	Id   uint64 `json:"id" gorm:"primary_key"`
	Uid  uint64 `json:"uid"`
	Fuid uint64 `json:"fuid"`
	CT   int64  `json:"c_t" gorm:"column:c_t"`
}

func AddFriend(uid, fuid uint64) *gorm.DB {
	//db.LogMode(true)
	d := &Friend{
		Uid:  uid,
		Fuid: fuid,
		CT:   time.Now().Unix(),
	}
	return db.Table("friends").Create(d)
}

/**
 * 先简单的获取用户的100个好友
 */
func FindFriends(uid uint64) ([]*Friend, error) {
	var list []*Friend
	res := db.Table("friends").Where("uid=?", uid).Limit(100).Find(&list)
	if res.Error != nil {
		return list, res.Error
	}
	return list, nil
}
