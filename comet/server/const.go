package server

import "time"

const ROOM_LIVE_TIME = 24 * 7 * time.Hour
const SESSION_LIVE_TIME = 3600 * 24 * time.Second

const (
	TYPE_COMMON_MSG  = float64(0)   //单个用户消息
	TYPE_ROOM_MSG    = float64(1)   //聊天室消息 到聊天室所有的人
	TYPE_JOIN_ROOM   = float64(2)   //进入房间
	TYPE_LEAVE_ROOM  = float64(3)   //退出房间
	TYPE_CREATE_ROOM = float64(4)   //创建房间
	TYPE_LOGIN       = float64(5)   //登录
	TYPE_LOGOUT      = float64(6)   //退出
	TYPE_NOTICE_MSG  = float64(10)  //通知类消息广播
	TYPE_REGISTER    = float64(11)  //注册设备生成deviceToken
	TYPE_TRANSPOND   = float64(12)  //将消息转发到其他系统微服务，并将其他系统获取的结果返回给客户端
	TYPE_PING        = float64(99)  //PING
	TYPE_PONG        = float64(100) //PONG
)
