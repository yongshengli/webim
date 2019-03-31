package models

import (
	"github.com/astaxie/beego"
	"webim/comet/common"
	"github.com/gomodule/redigo/redis"
	"errors"
	"log"
	"encoding/json"
	"net/rpc/jsonrpc"
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
	beego.Debug("session manager init")
}
/**
 在本机查找session
 */
func (m *sessionManager) GetSessionByUid(uid int) *Session {

	if _, ok := m.users[uid]; ok {
		return m.sessions[m.users[uid]]
	}
	return nil
}
func (m *sessionManager) AddSession(s *Session) bool{
	if len(s.DeviceToken)<1{
		return false
	}
	m.sessions[s.DeviceToken] = s
	if s.User.Id>0{
		m.users[s.User.Id] = s.DeviceToken
	}
	return true
}

func (m *sessionManager) DelSession(s *Session) bool{
	if len(s.DeviceToken)<1{
		return false
	}
	delete(m.sessions, s.DeviceToken)
	if s.RoomId!=""{
		room, _ := GetRoom(s.RoomId)
		if room!=nil{
			room.Leave(s.DeviceToken)
		}
	}
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
/**
 * 遍历所有机器给用户发消息
 */
func (m *sessionManager) SendMsg(deviceToken string, msg Msg) (bool, error){
	user, err := getDeviceTokenInfoByDeviceToken(deviceToken)
	if user == nil {
		return false, err
	}
	if ip, ok := user["IP"]; ok && len(ip)>1{
		addr := ip+""
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
	return false, errors.New("设备不在线")
}

func (m *sessionManager) Broadcast(msg Msg) (bool, error) {
	if msg.MsgType!=TYPE_BROADCAST_MSG{
		return false, errors.New("消息类型不是广播消息")
	}
	for _, session := range m.sessions {
		session.Send(&msg)
	}
	sMap, err := ServerManager.List()
	if err != nil {
		return false, err
	}
	for addr := range sMap {
		client, err := jsonrpc.Dial("tcp", addr)
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

func delSeviceTokenInfo(deviceToken string) (int, error){
	return common.RedisClient.Del([]string{deviceTokenKey(deviceToken)})
}
func saveDeviceTokenInfo(user User) (string, error){
	if len(user.DeviceToken)<1{
		return "", errors.New("DeviceToken为空")
	}
	jsonStr, err := json.Marshal(user)
	if err!=nil{
		beego.Error(err)
		return "", err
	}
	return common.RedisClient.Set(deviceTokenKey(user.DeviceToken), jsonStr, 3600*24)
}
func getDeviceTokenByUid(uid string) (string, error){
	return redis.String(common.RedisClient.Get(uidKey(uid)))
}
func getDeviceTokenInfoByDeviceToken(deviceToken string) (map[string]string, error){
	res := map[string]string{}
	replay, err := common.RedisClient.Get(deviceTokenKey(deviceToken))
	if replay==nil{
		return res, err
	}
	return redis.StringMap(replay, err)
}
func deviceTokenKey(deviceToken string) string{
	return "comet:token:" + deviceToken
}

func uidKey(uid string) string{
	return "comet:uid:"+uid
}