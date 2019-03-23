package models

import (
    "github.com/astaxie/beego"
    "fmt"
)

type Command struct {
    *Msg
    s *Session
}

func (c *Command) Run(s *Session) {
    c.s = s
    switch c.MsgType {
    case TYPE_CREATE_ROOM:
        c.createRoom()
    case TYPE_ROOM_MSG:
        c.roomMsg()
    case TYPE_JOIN_ROOM:
        c.joinRoom()
    case TYPE_LEAVE_ROOM:
        c.leaveRoom()
    }
}
func (c *Command) leaveRoom() {
    if _, ok := c.Data["room_id"]; !ok {
        fmt.Println("room_id 为空")
        return
    }
    roomId := c.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room != nil {
        room.Leave(RUser{SId: c.s.Id, User: *c.s.User, IP: c.s.IP})
    }
}
func (c *Command) joinRoom() {
    if _, ok := c.Data["room_id"]; !ok {
        beego.Warn("room_id 为空")
        return
    }
    roomId := c.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room == nil {
        data := make(map[string]interface{})
        data["content"] = "房间不存在"
        c.s.Send(NewMsg(TYPE_COMMON_MSG, data))
    } else {
        ru := RUser{SId: c.s.Id, User: *c.s.User, IP: c.s.IP}
        res, err := room.Join(RUser{SId: c.s.Id, User: *c.s.User, IP: c.s.IP})
        if err != nil {
            beego.Error(err)
        }
        if res {
            data := make(map[string]interface{})
            data["room_id"] = room.Id
            data["content"] = ru.User.Name + "进入房间"
            msg := NewMsg(TYPE_ROOM_MSG, data)
            room.Broadcast(msg)
        }
    }
}
func (c *Command) roomMsg() {
    if _, ok := c.Data["room_id"]; !ok {
        beego.Warn("room_id 为空")
        return
    }
    roomId := c.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room == nil {
        data := make(map[string]interface{})
        data["content"] = "房间不存在"
        c.s.Send(NewMsg(TYPE_COMMON_MSG, data))
    } else {
        room.Broadcast(c.Msg)
    }
}
func (c *Command) createRoom() {
    if _, ok := c.Msg.Data["room_id"]; !ok {
        beego.Warn("room_id 为空")
        return
    }
    roomId := c.Msg.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room == nil {
        NewRoom(roomId, "")
        data := make(map[string]interface{})
        data["content"] = "创建房间成功"
        c.s.Send(NewMsg(TYPE_COMMON_MSG, data))
    }
    m := *c.Msg
    m.MsgType = TYPE_JOIN_ROOM
    c.s.do(&m)
}
