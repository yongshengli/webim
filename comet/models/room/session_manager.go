package room

import (
	"github.com/astaxie/beego"
	"webim/comet/models"
	"errors"
	"webim/comet/common"
	"encoding/json"
)
var SessionManager *Manager

type Manager struct {
	users    map[int]string
	sessions map[string]*Session
}

func init() {
	SessionManager = &Manager{
		users:    make(map[int]string),
		sessions: make(map[string]*Session),
	}
	beego.Debug("manager init")
}
func (m *Manager) GetSessionByUid(uid int) *Session {

	if _, ok := m.users[uid]; ok {
		return m.sessions[m.users[uid]]
	}
	return nil
}
func (m *Manager) AddSession(s *Session) bool{
	m.sessions[s.Id] = s
	if s.User.Id>0{
		m.users[s.User.Id] = s.Id
	}
	return true
}

func (m *Manager) DelSession(s *Session) bool{
	delete(m.sessions, s.Id)
	if s.User.Id>0{
		delete(m.users, s.User.Id)
		return true
	}
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
		SessionNum: len(SessionManager.sessions),
	}
	return monitor
}

func (m *Manager) SendMsg(sId string, msg *models.Msg) (bool, error){
	if s, ok := m.sessions[sId]; ok{
		s.Send(msg)
		return true, nil
	}else{
		sessionJson, err := common.RedisClient.Get(sIdKey(sId))
		if err!=nil{
			return false, err
		}
		var sessMap map[string]string
		err = json.Unmarshal([]byte(sessionJson.(string)), &sessMap)
		if err!= nil{
			return false, err
		}

		return false, errors.New("没有找到用户"+sId)
	}
}
func (m *Manager) Broadcast(msg *models.Msg) bool {
	for _, session := range m.sessions {
		session.Send(msg)
	}
	return true
}
func (m *Manager) GetSessionIp(sId string){

}
func sIdKey(sId string) string{
	return "sessionId:" + sId
}