package room

import (
	"container/list"
	"github.com/astaxie/beego"
	"webim/comet/models"
)
var SessionManager *Manager

type Manager struct {
	rooms map[int]*Room
	users map[int]*Session
	sessionList *list.List
}

func init() {
	SessionManager = &Manager{
		rooms:       make(map[int]*Room),
		users:       make(map[int]*Session),
		sessionList: list.New(),
	}
	beego.Debug("manager init")
}
func (m *Manager) GetSessionByUid(uid int) *Session{

	if _, ok:=m.users[uid];ok{
		return m.users[uid]
	}
	return nil
}
func (m *Manager) AddSession(s *Session) bool{
	e := m.sessionList.PushBack(s)
	s.P = e
	if s.User.Id>0{
		m.users[s.User.Id] = s
	}
	return true
}

func (m *Manager) DelSession(s *Session) bool{
	if s.User.Id>0{
		p := m.users[s.User.Id].P
		m.sessionList.Remove(p)
		delete(m.users, s.User.Id)
		return true
	}
	for v:=m.sessionList.Front();s!=nil;v =v.Next(){
		if v.Value.(*Session) == s{
			m.sessionList.Remove(v)
		}
	}
	return true
}

func (m *Manager) GetRoom(roomId int) *Room{
	if _, ok := m.rooms[roomId]; ok {
		return m.rooms[roomId]
	}
	return nil
}

func (m *Manager) AddRoom(roomId int) bool{
	m.rooms[roomId] = NewRoom(roomId, m)
	return true
}

type Monitor struct {
	UserNum    int `json:"user_num"`
	SessionNum int `json:"conn_num"`
	RoomNum    int `json:"room_num"`
}

func Count() Monitor {
	monitor := Monitor{
		UserNum:    len(SessionManager.users),
		SessionNum: SessionManager.sessionList.Len(),
		RoomNum:    len(SessionManager.rooms),
	}
	return monitor
}

func (m *Manager) DelRoom(roomId int) bool {
	if _, ok := m.rooms[roomId]; ok {
		delete(m.rooms, roomId)
	}
	return true
}

func (m *Manager) Broadcast(msg *models.Msg) bool {
	for s := m.sessionList.Front(); s != nil; s = s.Next() {
		s.Value.(*Session).Send(msg)
	}
	return true
}