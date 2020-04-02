package models

import (
	"fmt"
	"time"
	"webim/comet/common"

	"github.com/jinzhu/gorm"
)

const DEVICE_TABLE_NUM = 1024

//Device 设备表结构
type Device struct {
	Id          uint   `json:"id" gorm:"primary_key"`
	DeviceId    string `json:"device_id"`
	UmengToken  string `json:"umeng_token"`
	MiToken     string `json:"mi_token"`
	HuaWeiToken string `json:"huawei_token"`
	VivoToken   string `json:"vivo_token"`
	OppoToken   string `json:"oppo_token"`
	Uid         uint   `json:"uid"`
	Username    string `json:"username"`
	LastActive  int64  `json:"last_active"`

	CT int64 `json:"c_t"`
	UT int64 `json:"u_t"`
}

//FindOneDevice 查找设备信息
func FindOneDevice(deviceId string) *Device {
	device := new(Device)
	deviceDb.Table(DeviceTableName(deviceId)).Where("device_id=?", deviceId).First(&device)
	return device
}

//InsertDevice 插入设备信息
func InsertDevice(deviceId string, data *Device) *gorm.DB {
	data.CT = time.Now().Unix()
	data.LastActive = data.CT
	data.UT = data.CT
	return deviceDb.Table(DeviceTableName(deviceId)).Create(data)
}

//UpdateDevice 更新设备信息
func UpdateDevice(deviceId string, data *Device) *gorm.DB {
	return deviceDb.Table(DeviceTableName(deviceId)).Where("device_id=?", deviceId).Update(data)
}

//UpdateLastActive 更新最后活跃时间
func UpdateLastActive(deviceId string) *gorm.DB {
	return deviceDb.Table(DeviceTableName(deviceId)).
		Where("device_id=?", deviceId).
		Update("last_active", time.Now().Unix())
}

//GetTableName 获取表名
func (d *Device) GetTableName(deviceId string) string {
	return DeviceTableName(deviceId)
}

//DeviceTableName 设备表名
func DeviceTableName(deviceId string) string {
	return fmt.Sprintf("device_%d", common.StrMod(deviceId, DEVICE_TABLE_NUM))
}
