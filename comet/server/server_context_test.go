package server

import (
    "testing"
    "fmt"
    "webim/comet/common"
)

func TestServerManger_Register(t *testing.T) {
    server := server{
        Host : "127.0.0.1",
        Port : "8000",
    }

    common.RedisInitTest()
    _, err := ServerManager.Register(server.Host, server)
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
    ServerManager.Remove(CurrentServer.Host)
    serJson, err := common.RedisClient.Do("hget", serverMapKey(), server.Host)
    if err != nil {
        t.Error(err)
    }
    if serJson != nil {
        t.Error("删除失败")
    }
}