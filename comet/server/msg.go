package server

import (
	"reflect"
	"errors"
	"fmt"
)

const (
	TYPE_COMMON_MSG    = float64(0) //单个用户消息
	TYPE_ROOM_MSG      = float64(1) //聊天室消息 到聊天室所有的人
	TYPE_JOIN_ROOM     = float64(2) //进入房间
	TYPE_LEAVE_ROOM    = float64(3) //退出房间
	TYPE_CREATE_ROOM   = float64(4) //创建房间
	TYPE_LOGIN         = float64(5)  //登录
	TYPE_LOGOUT        = float64(6)  //退出
	TYPE_NOTICE_MSG    = float64(10) //通知类消息
	TYPE_REGISTER      = float64(11) //注册设备生成deviceToken
	TYPE_TRANSPOND     = float64(12) //将消息转发到其他系统微服务，并将其他系统获取的结果返回给客户端
	TYPE_PING          = float64(99) //PING
	TYPE_PONG		   = float64(100) //PONG
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
	Type    float64                `json:"type" valid:"Required"`
	Version string                 `json:"version"`
	ReqId   string                 `json:"req_id"`
	Encode  string                 `json:"encode"`
	Data    map[string]interface{} `json:"data"`
}

func Map2Msg(data map[string]interface{}) (Msg, error){
	msg := &Msg{}
	elem := reflect.ValueOf(msg).Elem()
	relType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		tag := relType.Field(i).Tag
		mk := tag.Get("json")
		if mv, ok := data[mk]; ok {
			if elem.Field(i).Kind() != reflect.ValueOf(mv).Kind() {
				return *msg, errors.New(fmt.Sprintf("Map2Msg error: map[%s]type error", mk))
			}
			elem.Field(i).Set(reflect.ValueOf(mv))
		}
	}
	return *msg, nil
}