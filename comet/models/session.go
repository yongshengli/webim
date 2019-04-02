package models

import (
    "github.com/astaxie/beego"
    "github.com/gorilla/websocket"
    "encoding/json"
    "time"
    "io"
    "strings"
)

type User struct {
    Id            string `json:"id"`
    Name          string `json:"name"`
    Platform      string `json:"platform"`
    ClientVersion string `json:"clientVersion"`

    DeviceToken string                 `json:"device_token"` // CometToken = md5(udid+appKey)
    Info        map[string]interface{} `json:"info"`
    IP          string                 `json:"ip"`
    RealIP      string                 `json:"real_ip"`
}

type Session struct {
    DeviceToken string
    User        *User
    RoomId      string
    Conn        *websocket.Conn
    Manager     *sessionManager
    IP          string //用户所属机器ip
    reqChan     chan *Msg
    repChan     chan *Msg
    stopChan    chan bool
    sendFailCount int
}

func NewSession(conn *websocket.Conn, m *sessionManager) *Session {
    u := &User{
        Id:     "",
        Name:   "匿名用户",
        IP:     CurrentServer.Host,
        RealIP: conn.RemoteAddr().String(),
    }
    return &Session{
        DeviceToken: "",
        User:    u,
        Conn:    conn,
        Manager: m,
        IP:    CurrentServer.Host,
        stopChan: make(chan bool),
        reqChan: make(chan *Msg, 1000),
        repChan: make(chan *Msg, 1000),
        sendFailCount: 0,
    }
}

func (s *Session) Run() {
    defer s.Close()

    s.Manager.AddSession(s)
    go s.start()

    s.read()
}
func (s *Session) start() {
    ci, err := beego.AppConfig.Int64("heartbeat.interval")
    var interval time.Duration
    if err != nil {
        interval = time.Minute * 4
    } else {
        interval = time.Duration(ci)
    }
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <- s.stopChan:
            s.Close()
            return
        case req := <-s.reqChan:
            s.do(req)
        case rep := <-s.repChan:
            s.write(rep)

        case <-ticker.C:
            s.ping()
        }
    }
}
func (s *Session) Send(msg *Msg) {
    beego.Debug("session send call")
    s.repChan <- msg
}
func (s *Session)ping(){
    //当前session没有token信息则不发送ping,断开链接
    if s.checkSession() == false {
        s.Close()
        return
    }
    msg := &Msg{MsgType:TYPE_PING, Data:map[string]interface{}{"content":"ping"}}
    s.Send(msg)
}
//检查session是否有效
func (s *Session) checkSession() bool{
    return s.Manager.CheckSession(s)
}

func (s *Session) pong(){
    msg := &Msg{MsgType:TYPE_PONG, Data:map[string]interface{}{"content":"pong"}}
    s.Send(msg)
}
func (s *Session) write(msg *Msg) error {
    data, err := json.Marshal(msg)
    if err != nil {
        beego.Error("Fail to marshal event:", err)
        return err
    }
    err = s.Conn.WriteMessage(websocket.TextMessage, data)
    if err!= nil {
        s.sendFailCount ++
        if s.sendFailCount >= 3 {
            s.Close()
        }
        // 网络已经被关闭的情况下,设置Session关闭
        if err == io.EOF || err != nil && strings.Contains(err.Error(), "use of closed network connection") {
            beego.Info("msg[network_has_closed_than_set_session_close] sessionIp[%s] user[%v]", s.Conn.RemoteAddr(), s.User)
            s.sendFailCount = 9999
            s.Close()
        }
        return err
    }
    s.sendFailCount = 0
    return nil
}

func (s *Session) do(msg *Msg) {
    if msg.MsgType == TYPE_PONG {
        return
    } else if msg.MsgType == TYPE_PING {
        s.pong()
        return
    }
    //没有带deviceToken的链接不予许访问register以外的业务方法
    if msg.MsgType != TYPE_REGISTER{
        if s.checkSession() == false {
            if token, ok := msg.Data["device_token"]; ok{
                s.DeviceToken = token.(string)
                s.Manager.AddSession(s)
            }else{
                return
            }
        }
    }
    NewJobWork(*msg, s).Do()
}
func (s *Session) read() {
    for {
        if s.stopChan == nil {
            beego.Info("msg[stop_read_client_data] user[%v]", s.User)
            break
        }
        _, p, err := s.Conn.ReadMessage()
        if err != nil && err == io.EOF {
            beego.Warn("msg[disconnected_websocket] detail[%s]", err.Error())
            s.Close()
            break
        }
        if len(p)>0 {
            msg := new(Msg)
            json.Unmarshal(p, msg)
            s.reqChan <- msg
        }
    }
}
func (s *Session) Close() {
    defer s.Conn.Close()
    s.Manager.DelSession(s)

    if s.stopChan!=nil {
        close(s.stopChan)
        s.stopChan = nil
    }
}