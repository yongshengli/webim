package server

import (
	"comet/service/room"
	"runtime"

	"github.com/astaxie/beego/logs"
)

//StatusData 集群状态指标
type StatusData struct {
	UserNum      int `json:"user_num"`
	SessionNum   int `json:"conn_num"`
	RoomNum      int `json:"room_num"`
	ServerNum    int `json:"server_num"`
	GoroutineNum int `json:"goroutine_num"`
}

//Status 集群服务状态
func Status() StatusData {
	monitor := StatusData{
		UserNum:      0,
		SessionNum:   IMServer.CountSession(),
		GoroutineNum: runtime.NumGoroutine(),
	}
	roomNum, err := room.TotalRoom()
	if err != nil {
		logs.Error("msg[获取room num err] err[%s]", err.Error())
	}
	monitor.RoomNum = roomNum
	serverNum, err := IMServer.context.Len()
	if err != nil {
		logs.Error("msg[获取server num err] err[%s]", err.Error())
	}
	monitor.ServerNum = serverNum

	return monitor
}
