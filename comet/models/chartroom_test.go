package models

import (
    "testing"
    "github.com/satori/go.uuid"
)

func TestNewRoom(t *testing.T) {
    roomId := "123"
    room, err := NewRoom(roomId, "")
    if err!=nil{
        t.Error(err)
    }
    if room == nil{
        t.Fail()
    }
    room, err = GetRoom(roomId)
    if err!= nil{
       t.Error(err)
    }
    if room.Id =="" {
        t.Fail()
    }
    _, err = room.Join(RUser{SId:uuid.NewV4().String(), IP:"127.0.0.1", User:User{Id:1, Name:"张三"}})
    if err!=nil{
        t.Error(err)
    }
    users, err := room.Users()
    if err!= nil{
        t.Error(err)
    }
    if len(users)<1{
        t.Error("用户进入房间失败")
    }
    res, err :=  DelRoom(roomId)
    if err!=nil{
       t.Error(err)
    }
    if res<1{
       t.Error("删除聊天室失败")
    }
}