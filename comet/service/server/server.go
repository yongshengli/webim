package server

import (
	"comet/common"
	"comet/service/base"
	"comet/service/room"
	"comet/service/session"
	"errors"
	"fmt"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var IMServer *Server

//Info Server基本信息
type Info struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	LastActive int64  `json:"last_active"` //上次活跃时间
}
type Server struct {
	Info
	List             *[]Info
	context          *ServerContext
	users            sync.Map
	slotLen          int     //每个Solt 的长度
	slotContainerLen int     //session 容器长度
	slotContainer    []*Slot //session 容器
}

//NewServer NewServer
func NewServer(host, port string, slotContainerLen, slotLen int) *Server {
	info := Info{
		Host:       host,
		Port:       port,
		LastActive: time.Now().Unix(),
	}
	svrContext := new(ServerContext)
	s := &Server{
		Info:             info,
		List:             &[]Info{info},
		context:          svrContext,
		users:            sync.Map{},
		slotLen:          slotLen,
		slotContainerLen: slotContainerLen,
		slotContainer:    make([]*Slot, slotContainerLen),
	}
	svrContext.Server = s
	for i := 0; i < slotContainerLen; i++ {
		s.slotContainer[i] = NewSlot(slotLen)
	}
	s.context.Register(s.Info)
	return s
}

//Run 执行
func Run(host, port string, slotContainerLen, slotLen int) {
	if host == "" || port == "" {
		panic("host:port不能为空")
	}
	IMServer = NewServer(host, port, slotContainerLen, slotLen)
	go IMServer.ReportLive()
	go RunRpcService(IMServer)
	logs.Debug("msg[server start...]")
}

//UpdateList 更新server缓存列表
func (s *Server) UpdateList() {
	list, err := s.context.List()
	if err != nil {
		logs.Error("msg[更新server列表缓存失败] err[%s]", err.Error())
	}
	if len(list) > 0 {
		s.List = &list
	}
}

//ReportLive 服务报活
func (s *Server) ReportLive() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()
	for {
		<-t.C
		s.Info.LastActive = time.Now().Unix()
		s.context.Register(s.Info)
		s.UpdateList()
		logs.Debug("msg[服务报活,更新server缓存] server[%s:%s]", s.Host, s.Port)
	}
}
func (s *Server) getSlotPos(deviceToken string) int {
	return common.StrMod(deviceToken, s.slotContainerLen)
}

func (s *Server) getSlot(deviceToken string) *Slot {
	return s.slotContainer[s.getSlotPos(deviceToken)]
}

//CheckSession 检查session用户是否登录
func (s *Server) CheckSession(ss *session.Session) bool {
	if len(ss.DeviceToken) < 1 {
		return false
	}
	return s.getSlot(ss.DeviceToken).Has(ss.DeviceToken)
}

//GetSessionByDeviceToken 根据tocken查找session
func (s *Server) GetSessionByDeviceToken(deviceToken string) *session.Session {
	return s.getSlot(deviceToken).Get(deviceToken)
}

//GetSessionByUid 在本机查找session
func (s *Server) GetSessionByUid(uid string) *session.Session {
	if deviceToken, ok := s.users.Load(uid); ok {
		t := deviceToken.(string)
		return s.getSlot(t).Get(t)
	}
	return nil
}

//CountSession 统计本机session的数量
func (s *Server) CountSession() int {
	num := 0
	for i := 0; i < s.slotContainerLen; i++ {
		num += s.slotContainer[i].Len()
	}
	return num
}

//AddSession AddSession
func (s *Server) AddSession(ss *session.Session) bool {
	if len(ss.DeviceToken) < 1 {
		return false
	}
	ss.User.DeviceToken = ss.DeviceToken
	//把session 信息保存到redis
	_, err := saveDeviceTokenInfo(ss.User)
	if err != nil {
		logs.Error("msg[AddSession err] err[%s]", err.Error())
		return false
	}
	s.getSlot(ss.DeviceToken).Add(ss.DeviceToken, ss)
	if ss.User.Id != 0 {
		s.users.Store(ss.User.Id, ss.DeviceToken)
	}
	return true
}

//DelSession 删除Session
func (s *Server) DelSession(ss *session.Session) bool {
	if len(ss.DeviceToken) < 1 {
		return false
	}
	//从redis中删除session
	_, err := delDeviceTokenInfo(ss.DeviceToken)
	if err != nil {
		beego.Error(err)
	}
	s.getSlot(ss.DeviceToken).Del(ss.DeviceToken)
	if ss.RoomId != "" {
		myRoom, _ := room.GetRoom(ss.RoomId)
		if myRoom != nil {
			myRoom.Leave(ss)
		}
	}
	if ss.User.Id != 0 {
		s.users.Delete(ss.User.Id)
		return true
	}
	return true
}

//Unicast 根据deviceToken找到用户对应的主机然后推送给用户
func (s *Server) Unicast(deviceToken string, msg base.Msg) (bool, error) {
	user, err := getDeviceTokenInfo(deviceToken)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("用户token不存在")
	}
	if len(user.IP) > 1 {
		if user.IP == fmt.Sprintf("%s:%s", s.Host, s.Port) {
			return s.SendMsg(deviceToken, msg)
		}
		client, err := jsonrpc.Dial("tcp", user.IP)
		if err != nil {
			logs.Error("msg[连接Dial的发生了错误addr:%s], err[%s]", user.IP, err.Error())
			return false, err
		}
		defer client.Close()
		args := map[string]interface{}{}
		args["device_token"] = deviceToken
		args["msg"] = msg
		reply := false
		err = client.Call("RpcService.Unicast", args, &reply)
		if err != nil {
			logs.Error("msg[发送单播addr:%s失败] args[%v] res[%v] err[%s]", user.IP, args, reply, err.Error())
		} else {
			logs.Debug("msg[发送单播addr:%s成功] args[%v] res[%v]", user.IP, args, reply)
		}
		return true, nil
	}
	return false, errors.New("设备不在线")
}

//SendMsg 向本机用户发送消息
func (s *Server) SendMsg(deviceToken string, msg base.Msg) (bool, error) {
	logs.Debug("msg[call_SendMsg] device_token[%s]", deviceToken)
	slot := s.getSlot(deviceToken)
	if slot.Has(deviceToken) {
		slot.Get(deviceToken).Send(msg)
		return true, nil
	}
	delDeviceTokenInfo(deviceToken)
	return false, errors.New("设备不在线")
}

//Broadcast 全部在线用户消息广播
func (s *Server) Broadcast(msg base.Msg) (bool, error) {
	logs.Debug("msg[call_Broadcast]")
	for _, st := range *s.List {
		if st.Host == s.Host {
			s.BroadcastSelf(msg)
		} else {
			addr := st.Host + ":" + st.Port
			client, err := jsonrpc.Dial("tcp", addr)
			if err != nil {
				logs.Error("msg[连接Dial的发生了错误] addr[%s], err:%s", addr, err.Error())
				continue
			}
			args := map[string]interface{}{}
			args["msg"] = msg
			reply := false
			client.Call("RpcService.BroadcastSelf", args, &reply)
			client.Close()
			logs.Debug("msg[发送广播] addr[%s], res:%t", addr, reply)
		}
	}
	return true, nil
}

//BroadcastSelf 本机在线用户消息广播
func (s *Server) BroadcastSelf(msg base.Msg) (bool, error) {
	logs.Debug("msg[call_BroadcastSelf]")
	for _, slot := range s.slotContainer {
		sessionsMap := slot.All()
		for _, session := range sessionsMap {
			session.Send(msg)
		}
	}

	return true, nil
}
