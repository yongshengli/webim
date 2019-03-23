package models

import (
	"github.com/astaxie/beego"
	"errors"
	"net/rpc"
	"log"
	"webim/comet/common"
)
var SessionManager *sessionManager

type sessionManager struct {
	users    map[int]string
	sessions map[string]*Session
}

func init() {
	SessionManager = &sessionManager{
		users:    make(map[int]string),
		sessions: make(map[string]*Session),
	}
	beego.Debug("manager init")
}
func (m *sessionManager) GetSessionByUid(uid int) *Session {

	if _, ok := m.users[uid]; ok {
		return m.sessions[m.users[uid]]
	}
	return nil
}
func (m *sessionManager) AddSession(s *Session) bool{
	m.sessions[s.Id] = s
	if s.User.Id>0{
		m.users[s.User.Id] = s.Id
	}
	return true
}

func (m *sessionManager) DelSession(s *Session) bool{
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

func (m *sessionManager) SendMsgAll(sId string, msg Msg) (bool, error){
	if _, ok := m.sessions[sId]; ok{
		return m.SendMsg(sId, msg)
	}else{
		sMap, err := ServerManager.List()
		if err != nil {
			return false, err
		}
		localAddr := common.GetLocalIp()
		for ip, port := range sMap{
			addr := ip + ":" + port
			if localAddr == addr {
				continue
			}
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
			log.Printf("发送广播addr%s, res:%t", addr, reply)
		}
		return true, nil
	}
}

func (m *sessionManager) SendMsg(sId string, msg Msg) (bool, error){
	if s, ok := m.sessions[sId]; ok {
		s.Send(&msg)
		return true, nil
	}
	return false, errors.New("没有找到用户"+sId)
}

func (m *sessionManager) Broadcast(msg Msg) (bool, error) {
	for _, session := range m.sessions {
		session.Send(&msg)
	}
	sMap, err := ServerManager.List()
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
		log.Printf("发送广播addr%s, res:%t", addr, reply)
	}
	return true, nil
}

func (m *sessionManager) BroadcastSelf(msg Msg) (bool, error) {
	for _, session := range m.sessions {
		session.Send(&msg)
	}
	return true, nil
}
func (m *sessionManager) GetSessionIp(sId string){

}
func sIdKey(sId string) string{
	return "comet:sessionId:" + sId
}