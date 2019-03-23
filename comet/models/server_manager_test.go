package models

import (
    "testing"
    "fmt"
    "webim/comet/common"
)

func TestServerManger_Register(t *testing.T) {
    common.RedisInitTest()
    _, err := ServerManager.Register("8000")
    if err!=nil{
        t.Error(err)
    }
    //fmt.Println(res)
    sMap, err := ServerManager.List()
    if err!= nil{
        t.Error(err)
    }
    if len(sMap)<1{
        t.Error("没有取到主机map")
    }
    fmt.Println(sMap)
    ServerManager.Remove()
}