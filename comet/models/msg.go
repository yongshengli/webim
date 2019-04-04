package models

import (
	"reflect"
)

const (
	TYPE_COMMON_MSG    = 0 //单个用户消息
	TYPE_ROOM_MSG      = 1 //聊天室消息 到聊天室所有的人
	TYPE_JOIN_ROOM     = 2 //进入房间
	TYPE_LEAVE_ROOM    = 3 //退出房间
	TYPE_CREATE_ROOM   = 4 //创建房间
	TYPE_LOGIN         = 5  //登录
	TYPE_LOGOUT        = 6  //退出
	TYPE_NOTICE_MSG    = 10 //通知类消息
	TYPE_REGISTER      = 11 //注册设备生成deviceToken
	TYPE_TRANSPOND     = 12 //将消息转发到其他系统微服务，并将其他系统获取的结果返回给客户端
	TYPE_PING          = 99 //PING
	TYPE_PONG		   = 100 //PONG
)

type Job struct {
	TraceID string // 日志ID

	ReqTime int64 // 请求到达时间
	RspTime int64 // 响应结束时间

	TransferReqTime int64 // 转发请求开始时间
	TransferRspTime int64 // 转发请求结束时间

	Req Msg // 原请求信息
	Rsp Msg // 响应信息
}

type Msg struct {
	Type    int                    `json:"type" valid:"Required"`
	Version string                 `json:"version"`
	ReqId   string                 `json:"req_id"`
	Encode  string                 `json:"encode"`
	Data    map[string]interface{} `json:"data"`
}

func Map2Msg(m map[string]interface{}) Msg{
	msg := &Msg{}

	elem := reflect.ValueOf(msg).Elem()
	relType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		tag := relType.Field(i).Tag
		mK := tag.Get("json")
		if _, ok := m[mK]; ok {
			if elem.Field(i).Type() == reflect.ValueOf(m[mK]).Type() {
				elem.Field(i).Set(reflect.ValueOf(m[mK]))
			}
		}
	}
	return *msg
}