package models

import (
    "github.com/astaxie/beego"
    "github.com/gorilla/websocket"
    "encoding/json"
    "github.com/satori/go.uuid"
)

type User struct {
    Id            int                    `json:"id"`
    Name          string                 `json:"name"`
    Platform      string                 `json:"platform"`
    ClientVersion string                 `json:"clientVersion"`

    DeviceToken    string                 `json:"DeviceToken"` // CometToken = md5(udid+appKey)
    Info          map[string]interface{} `json:"info"`
    IP            string                 `json:"ip"`
}

type Session struct {
    Id      string
    User    *User
    RoomId  string
    Conn    *websocket.Conn
    Manager *sessionManager
    Addr    string //用户所属机器ip
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
        Addr:    CurrentServer["host"] + ":" + CurrentServer["port"],
        reqChan: make(chan *Msg, 1000),
        repChan: make(chan *Msg, 1000),
    }
}

func (s *Session) Run() {
    defer s.Close()

    s.Manager.AddSession(s)
    go s.start()
    s.read()
}
func (s *Session) start() {
    for {
        select {
        case req := <-s.reqChan:
            s.do(req)
        case rep := <-s.repChan:
            s.write(rep)
        }
    }
}
func (s *Session) Send(msg *Msg) {
    beego.Debug("session send call")
    s.repChan <- msg
}

func (s *Session) write(msg *Msg) bool {
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

func (s *Session) do(msg *Msg) {
    NewJobWork(*msg, s).Do()
}
func (s *Session) read() {
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
func (s *Session) Close() {
    defer s.Conn.Close()
    s.Manager.DelSession(s)
}