package room

import (
	"container/list"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"encoding/json"
	"webim/comet/models"
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
	Manager *Manager
	P *list.Element
	Ip string //用户所属机器ip
	reqChan chan *models.Msg
	repChan chan *models.Msg
}
func NewSession(conn *websocket.Conn, m *Manager) *Session{
	u := &User{
		Id:0,
		Name:"匿名用户",
	}
	return &Session{
		User: u,
		Conn:conn, Manager:m,
		reqChan:make(chan *models.Msg, 1000),
		repChan:make(chan *models.Msg, 1000),
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
func (s *Session) Send(msg *models.Msg){
	beego.Debug("session send call")
	s.repChan <- msg
}

func (s * Session) write(msg *models.Msg) bool{
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

func (s *Session) do(msg *models.Msg){
	switch msg.MsgType {
	case models.TYPE_CREATE_ROOM:
		if _, ok := msg.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := int(msg.Data["room_id"].(float64))
		room := s.Manager.GetRoom(roomId)
		if room == nil {
			s.Manager.AddRoom(roomId)
			data := make(map[string]interface{})
			data["content"] = "创建房间成功"
			s.Send(models.NewMsg(models.TYPE_COMMON_MSG, data))
		}
		m := *msg
		m.MsgType = models.TYPE_JOIN_ROOM
		s.do(&m)
	case models.TYPE_ROOM_MSG:
		if _, ok := msg.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := int(msg.Data["room_id"].(float64))
		room := s.Manager.GetRoom(roomId)
		if room == nil {
			data := make(map[string]interface{})
			data["content"] = "房间不存在"
			s.Send(models.NewMsg(models.TYPE_COMMON_MSG, data))
		} else {
			room.Broadcast(msg)
		}
	case models.TYPE_JOIN_ROOM:
		if _, ok := msg.Data["room_id"]; !ok {
			beego.Warn("room_id 为空")
			return
		}
		roomId := int(msg.Data["room_id"].(float64))
		room := s.Manager.GetRoom(roomId)
		if room == nil {
			data := make(map[string]interface{})
			data["content"] = "房间不存在"
			s.Send(models.NewMsg(models.TYPE_COMMON_MSG, data))
		} else {
			room.Join(*s)
		}
	case models.TYPE_LEAVE_ROOM:
		if _, ok := msg.Data["room_id"]; !ok {
			fmt.Println("room_id 为空")
			return
		}
		roomId := int(msg.Data["room_id"].(float64))
		room := s.Manager.GetRoom(roomId)
		if room != nil {
			room.Leave(*s)
		}
	}
}
func (s *Session) read(){
	for {
		_, p, err := s.Conn.ReadMessage()
		if err != nil {
			return
		}
		msg := new(models.Msg)
		json.Unmarshal(p, msg)
		s.reqChan <- msg
	}
}
func (s *Session) Close(){
	s.Conn.Close()
	s.Manager.DelSession(s)
}