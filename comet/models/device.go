package models

import (
    "webim/comet/common"
    "fmt"
    "time"
    "github.com/jinzhu/gorm"
)

const DEVICE_TABLE_NUM = 1024

type Device struct {
    Id          uint   `json:"id" gorm:"primary_key"`
    UmengToken  string `json:"umeng_token"`
    MiToken     string `json:"mi_token"`
    HuaWeiToken string `json:"huawei_token"`
    VivoToken   string `json:"vivo_token"`
    OppoToken   string `json:"oppo_token"`
    Uid         uint   `json:"uid"`
    LastActive  int64  `json:"last_active"`

    CT          int64  `json:"c_t"`
    UT          int64  `json:"u_t"`
}

func FindOneDevice(deviceToken string) *Device {
    device := new(Device)
    deviceDb.Table(DeviceTableName(deviceToken)).Where("umeng_token=?", deviceToken).First(&device)
    return device
}
//插入设备信息
func InsertDevice(deviceToken string, data *Device) *gorm.DB {
    data.CT = time.Now().Unix()
    data.LastActive = data.CT
    data.UT = data.CT
    return deviceDb.Table(DeviceTableName(deviceToken)).Create(data)
}
//更新设备信息
func UpdateDevice(deviceToken string, data *Device) *gorm.DB {
    return deviceDb.Table(DeviceTableName(deviceToken)).Where("umeng_token=?", deviceToken).Update(data)
}
//更新最后活跃时间
func UpdateLastActive(deviceToken string) *gorm.DB {
    return deviceDb.Table(DeviceTableName(deviceToken)).
        Where("umeng_token=?", deviceToken).
        Update("last_active", time.Now().Unix())
}

func (d *Device) GetTableName(deviceToken string) string {
    return DeviceTableName(deviceToken)
}

func DeviceTableName(deviceToken string) string {
    return fmt.Sprintf("device_%d", common.StrMod(deviceToken, DEVICE_TABLE_NUM))
}