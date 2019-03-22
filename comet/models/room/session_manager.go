package room

import (
	"github.com/astaxie/beego"
	"webim/comet/models"
	"errors"
	"webim/comet/common"
	"encoding/json"
	"net/rpc"
	"log"
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

func (m *Manager) SendMsgAll(sId string, msg models.Msg) (bool, error){
	if _, ok := m.sessions[sId]; ok{
		return m.SendMsg(sId, msg)
	}else{
		sessionJson, err := common.RedisClient.Get(sIdKey(sId))
		if err!=nil{
			return false, err
		}
		var sessMap map[string]string
		err = json.Unmarshal(sessionJson.([]byte), &sessMap)
		if err!= nil{
			return false, err
		}
		if sessMap["IP"]==common.GetLocalIp(){
			return false, errors.New("没有找到用户"+sId)
		}
		sMap, err := models.ServerManager.List()
		if err != nil {
			return false, err
		}
		for ip, port := range sMap{
			addr := ip + ":" + port
			client, err := rpc.Dial("tcp", addr)
			if err != nil {
				log.Printf("连接Dial的发生了错误addr:%s, err:%s", addr, err.Error())
				continue
			}
			args := map[string]interface{}{}
			args["sid"] = sId
			args["msg"] = msg
			reply := false
			client.Call("RpcFunc.Unicast", args, &reply)
			log.Printf("发送广播addr%s, res:%s", addr, reply)
		}
		return true, nil
	}
}

func (m *Manager) SendMsg(sId string, msg models.Msg) (bool, error){
	if s, ok := m.sessions[sId]; ok {
		s.Send(&msg)
		return true, nil
	}
	return false, errors.New("没有找到用户"+sId)
}

func (m *Manager) Broadcast(msg models.Msg) (bool, error) {
	for _, session := range m.sessions {
		session.Send(&msg)
	}
	sMap, err := models.ServerManager.List()
	if err != nil {
		return false, err
	}
	for ip, port := range sMap {
		addr := ip + ":" + port
		client, err := rpc.Dial("tcp", addr)
		if err != nil {
			log.Printf("连接Dial的发生了错误addr:%s, err:%s", addr, err.Error())
			continue
		}
		args := map[string]interface{}{}
		args["msg"] = msg
		reply := false
		client.Call("RpcFunc.Broadcast", args, &reply)
		log.Printf("发送广播addr%s, res:%s", addr, reply)
	}
	return true, nil
	return true, nil
}

func (m *Manager) BroadcastSelf(msg models.Msg) (bool, error) {
	for _, session := range m.sessions {
		session.Send(&msg)
	}
	return true, nil
}
func (m *Manager) GetSessionIp(sId string){

}
func sIdKey(sId string) string{
	return "comet:sessionId:" + sId
}