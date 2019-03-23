package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type User struct {
	Id int
	Name string
	Info map[string]interface{}
}

type Session struct {
	Id string
	User *User
	Conn  *websocket.Conn
	Manager *sessionManager
	IP string //用户所属机器ip
	reqChan chan *Msg
	repChan chan *Msg
}

func NewSession(conn *websocket.Conn, m *sessionManager) *Session {
	u := &User{
		Id:   0,
		Name: "匿名用户",
	}
	return &Session{
		Id:      uuid.NewV4().String(),
		User:    u,
		Conn:    conn,
		Manager: m,
		IP:      conn.LocalAddr().String(),
		reqChan: make(chan *Msg, 1000),
		repChan: make(chan *Msg, 1000),
	}
}

func (s *Session) Run(){
	defer s.Close()

	s.Manager.AddSession(s)
	go s.start()
	s.read()
}
func (s *Session) start(){
	for {
		select {
		case req := <-s.reqChan:
			s.do(req)
		case rep := <-s.repChan:
			s.write(rep)
		}
	}
}
func (s *Session) Send(msg *Msg){
	beego.Debug("session send call")
	s.repChan <- msg
}

func (s * Session) write(msg *Msg) bool{
	data, err := json.Marshal(msg)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return false
	}
	if s.Conn.WriteMessage(websocket.TextMessage, data) != nil {
		// User disconnected. delete from room
		s.Close()
	}
	return true
}

func (s *Session) do(msg *Msg){
	switch msg.MsgType {
	case TYPE_CREATE_ROOM:
		if _, ok := msg.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := msg.Data["room_id"].(string)
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
		m := *msg
		m.MsgType = TYPE_JOIN_ROOM
		s.do(&m)
	case TYPE_ROOM_MSG:
		if _, ok := msg.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := msg.Data["room_id"].(string)
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
			room.Broadcast(msg)
		}
	case TYPE_JOIN_ROOM:
		if _, ok := msg.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := msg.Data["room_id"].(string)
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
		if _, ok := msg.Data["room_id"]; !ok {
			fmt.Println("room_id 为空")
			return
		}
		roomId := msg.Data["room_id"].(string)
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
func (s *Session) read(){
	for {
		_, p, err := s.Conn.ReadMessage()
		if err != nil {
			return
		}
		msg := new(Msg)
		json.Unmarshal(p, msg)
		s.reqChan <- msg
	}
}
func (s *Session) Close(){
	s.Conn.Close()
	s.Manager.DelSession(s)
}