package models

const (
	TYPE_COMMON_MSG    = 0 //单个用户消息
	TYPE_ROOM_MSG      = 1 //聊天室消息 到聊天室所有的人
	TYPE_JOIN_ROOM     = 2 //进入房间
	TYPE_LEAVE_ROOM    = 3 //退出房间
	TYPE_CREATE_ROOM   = 4 //创建房间
	TYPE_LOGIN         = 5  //登录
	TYPE_LOGOUT        = 6  //退出
	TYPE_BROADCAST_MSG = 10 //广播消息 到全部用户
	TYPE_REGISTER      = 11 //注册设备生成deviceToken
	TYPE_TRANSPOND     = 12 //将消息转发到其他系统微服务，并将其他系统获取的结果返回给客户端
)
type MsgType int

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
	MsgType MsgType                `json:"type"`
	Version string                 `json:"version"`
	ReqId   string                 `json:"req_id"`
	Encode  string                 `json:"encode"`
	Data    map[string]interface{} `json:"data"`
}