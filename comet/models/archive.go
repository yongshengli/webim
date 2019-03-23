package models

import (
	"github.com/astaxie/beego"
	"fmt"
)

type Command struct {
	*Msg
}

func (c *Command) Run(s *Session){
	switch c.MsgType {
	case TYPE_CREATE_ROOM:
		if _, ok := c.Msg.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := c.Msg.Data["room_id"].(string)
		room, err := GetRoom(roomId)
		if err!=nil{
			beego.Error(err)
			return
		}
		if room == nil {
			NewRoom(roomId, "")
			data := make(map[string]interface{})
			data["content"] = "创建房间成功"
			s.Send(NewMsg(TYPE_COMMON_MSG, data))
		}
		m := *c.Msg
		m.MsgType = TYPE_JOIN_ROOM
		s.do(&m)
	case TYPE_ROOM_MSG:
		if _, ok := c.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := c.Data["room_id"].(string)
		room, err := GetRoom(roomId)
		if err!=nil{
			beego.Error(err)
			return
		}
		if room == nil {
			data := make(map[string]interface{})
			data["content"] = "房间不存在"
			s.Send(NewMsg(TYPE_COMMON_MSG, data))
		} else {
			room.Broadcast(c.Msg)
		}
	case TYPE_JOIN_ROOM:
		if _, ok := c.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := c.Data["room_id"].(string)
		room, err := GetRoom(roomId)
		if err!=nil{
			beego.Error(err)
			return
		}
		if room == nil {
			data := make(map[string]interface{})
			data["content"] = "房间不存在"
			s.Send(NewMsg(TYPE_COMMON_MSG, data))
		} else {
			ru := RUser{SId:s.Id,User:*s.User, IP:s.IP}
			res, err := room.Join(RUser{SId:s.Id,User:*s.User, IP:s.IP})
			if err!=nil{
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
	case TYPE_LEAVE_ROOM:
		if _, ok := c.Data["room_id"]; !ok {
			fmt.Println("room_id 为空")
			return
		}
		roomId := c.Data["room_id"].(string)
		room , err := GetRoom(roomId)
		if err != nil {
			beego.Error(err)
			return
		}
		if room != nil {
			room.Leave(RUser{SId:s.Id,User:*s.User, IP:s.IP})
		}
	}
}
