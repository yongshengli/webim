package common

import (
    "testing"
    "fmt"
)

func TestGetLocalIp(t *testing.T) {
    ip := GetLocalIp()
    if ip==""{
        t.Error("获取本机ip错误")
    }
    fmt.Println("本机ip: "+ip)
}
