package models

import "testing"

func TestServerManger_Register(t *testing.T) {
    res, err := ServerManager.Register("8000")
    if err!=nil{
        t.Error(err)
    }
    if res<1{
        t.Error("注册主机失败")
    }
    sMap, err := ServerManager.List()
    if err!= nil{
        t.Error(err)
    }
    if len(sMap)<1{
        t.Error("没有取到主机map")
    }
    ServerManager.Remove()
}