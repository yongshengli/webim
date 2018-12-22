package models

import "time"

const (
	TYPE_COMMON_MSG = 0 //普通长连接消息
	TYPE_ROOM_MSG   = 1 //聊天室消息
	TYPE_JOIN_ROOM  = 2 //进入房间
	TYPE_LEAVE_ROOM = 3 //退出房间
	TYPE_LOGIN      = 4 //登录
	TYPE_LOGOUT     = 5 //登出
)
type MsgType int

type Msg struct {
	MsgType MsgType                `json:"type"`
	Data    map[string]interface{} `json:"data"`
	Time    int64                    `json:"time"`
}
func New(t MsgType, d map[string]interface{}) *Msg{
	return &Msg{
		MsgType:t,
		Data: d,
		Time: time.Now().Unix(),
	}
}