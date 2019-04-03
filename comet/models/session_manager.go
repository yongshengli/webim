package models

import (
	"github.com/astaxie/beego"
	"webim/comet/common"
	"github.com/gomodule/redigo/redis"
	"errors"
	"log"
	"encoding/json"
	"net/rpc/jsonrpc"
    "time"
	"github.com/astaxie/beego/logs"
)
var SessionManager *sessionManager

type sessionManager struct {
	users    map[string]string
	sessions map[string]*Session
}

func init() {
	SessionManager = &sessionManager{
		users:    make(map[string]string),
		sessions: make(map[string]*Session),
	}
	beego.Debug("session manager init")
}

func (m *sessionManager) CheckSession(s *Session) bool{
	if len(s.DeviceToken) < 1 {
		return false
	}
	if _, ok := m.sessions[s.DeviceToken]; ok{
		return true
	}
	return false
}

/**
 在本机查找session
 */
func (m *sessionManager) GetSessionByUid(uid string) *Session {
	if _, ok := m.users[uid]; ok {
		return m.sessions[m.users[uid]]
	}
	return nil
}
func (m *sessionManager) AddSession(s *Session) bool{
	if len(s.DeviceToken)<1{
		return false
	}
	s.User.DeviceToken = s.DeviceToken
	//把session 信息保存到redis
	_, err := saveDeviceTokenInfo(s.User)
	if err!=nil{
		beego.Error(err)
		return false
	}
	m.sessions[s.DeviceToken] = s
	if s.User.Id!=""{
		m.users[s.User.Id] = s.DeviceToken
	}
	return true
}

func (m *sessionManager) DelSession(s *Session) bool{
	if len(s.DeviceToken)<1{
		return false
	}
	//从redis中删除session
	_, err := delDeviceTokenInfo(s.DeviceToken)
	if err!= nil{
		beego.Error(err)
	}
	delete(m.sessions, s.DeviceToken)
	if s.RoomId != "" {
		room, _ := GetRoom(s.RoomId)
		if room != nil {
			room.Leave(s)
		}
	}
	if s.User.Id != "" {
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
/**
 * 根据deviceToken找到用户对应的主机然后推送给用户
 */
func (m *sessionManager) Unicast(deviceToken string, msg Msg) (bool, error){
	user, err := getDeviceTokenInfoByDeviceToken(deviceToken)
	if user == nil {
		return false, err
	}
	if ip, ok := user["ip"]; ok && len(ip)>1{
		if ip == CurrentServer.Host {
			return m.SendMsg(deviceToken, msg)
		} else {
			addr := ip + ":" + CurrentServer.Port
			client, err := jsonrpc.Dial("tcp", addr)
			if err != nil {
				beego.Error("连接Dial的发生了错误addr:%s, err:%s", addr, err.Error())
				return false, err
			}
			args := map[string]interface{}{}
			args["device_token"] = deviceToken
			args["msg"] = msg
			reply := false
			client.Call("RpcFunc.Unicast", args, &reply)
			log.Printf("发送单播addr%s, res:%t", addr, reply)
			return true, nil
		}
	}
	return false, errors.New("设备不在线")
}
func (m *sessionManager) SendMsg(deviceToken string, msg Msg) (bool, error){
	logs.Debug("msg[call_SendMsg] device_token[%s]", deviceToken)
	if s, ok :=m.sessions[deviceToken]; ok{
		s.Send(&msg)
		return true, nil
	}else{
		delDeviceTokenInfo(deviceToken)
		return false, errors.New("设备不在线")
	}
}
func (m *sessionManager) Broadcast(msg Msg) (bool, error) {
	for _, session := range m.sessions {
		session.Send(&msg)
	}
	sMap, err := ServerManager.List()
	if err != nil {
		return false, err
	}
	for _, st := range sMap {
		if st.Host == CurrentServer.Host {
			return m.BroadcastSelf(msg)
		}else {
			addr := st.Host + ":" + st.Port
			client, err := jsonrpc.Dial("tcp", addr)
			if err != nil {
				logs.Error("msg[连接Dial的发生了错误] addr[%s], err:%s", addr, err.Error())
				continue
			}
			args := map[string]interface{}{}
			args["msg"] = msg
			reply := false
			client.Call("RpcFunc.BroadcastSelf", args, &reply)
			logs.Debug("msg[发送广播] addr[%s], res:%t", addr, reply)
		}
	}
	return true, nil
}

func (m *sessionManager) BroadcastSelf(msg Msg) (bool, error) {
	for _, session := range m.sessions {
		session.Send(&msg)
	}
	return true, nil
}

func delDeviceTokenInfo(deviceToken string) (int, error){
	logs.Debug("msg[call_delDeviceTokenInfo] device_toke[%s]", deviceToken)
	return common.RedisClient.Del([]string{deviceTokenKey(deviceToken)})
}
func saveDeviceTokenInfo(user *User) (string, error){
	if len(user.DeviceToken)<1{
		return "", errors.New("DeviceToken为空")
	}
	jsonStr, err := json.Marshal(user)
	if err!=nil{
		beego.Error(err)
		return "", err
	}
	return common.RedisClient.Set(deviceTokenKey(user.DeviceToken), jsonStr, 3600*24*time.Second)
}
func getDeviceTokenByUid(uid string) (string, error){
	return redis.String(common.RedisClient.Get(uidKey(uid)))
}

func getDeviceTokenInfoByDeviceToken(deviceToken string) (map[string]string, error){
	var res map[string]string
	replay, err := common.RedisClient.Get(deviceTokenKey(deviceToken))
	if replay == nil {
		return nil, err
	}
	err = json.Unmarshal(replay.([]byte), &res)
	if err != nil {
		logs.Error("msg[解析json失败] method[getDeviceTokenInfoByDeviceToken] err[%s]", err.Error())
		return nil, err
	}
	return res, nil
}
func deviceTokenKey(deviceToken string) string{
	return "comet:token:" + deviceToken
}

func uidKey(uid string) string{
	return "comet:uid:"+uid
}