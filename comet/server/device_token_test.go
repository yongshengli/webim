package server

import (
    "testing"
    "webim/comet/common"
    "fmt"
)

func TestDeviceTokenInfo(t *testing.T){
    common.RedisInitTest()

    u := &User{Id:"zhangsan", Name:"张三", IP:"127.0.0.1",DeviceToken:"1111111", DeviceId:"123321"}
    _, err := saveDeviceTokenInfo(u)
    if err != nil {
        t.Error(err)
    }
    uInfo, err := getDeviceTokenInfo(u.DeviceToken)
    if err!= nil{
        t.Error(err)
    }
    fmt.Println(uInfo)
    if v, ok := uInfo["id"]; !ok || v != u.Id {
        t.Errorf("获取的用户token信息错误id 期望 %s get %s", u.Id, v)
    }
    res, err := delDeviceTokenInfo(u.DeviceToken)

    if err != nil {
        t.Error(err)
    }
    if res < 1 {
        t.Error("删除devive_token失败")
    }
    uInfo, err = getDeviceTokenInfo(u.DeviceToken)
    if err!= nil{
        t.Error(err)
    }
}
