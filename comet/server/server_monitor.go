package server

import "github.com/astaxie/beego/logs"

type StatusData struct {
    UserNum    int `json:"user_num"`
    SessionNum int `json:"conn_num"`
    RoomNum    int `json:"room_num"`
    ServerNum  int `json:"server_num"`
}

func Status() StatusData {
    monitor := StatusData{
        UserNum:    0,
        SessionNum: Server.CountSession(),
    }
    roomNum, err := TotalRoom()
    if err != nil {
        logs.Error("msg[获取room num err] err[%s]", err.Error())
    }
    monitor.RoomNum = roomNum
    serverNum, err := Server.context.Len()
    if err!=nil {
        logs.Error("msg[获取server num err] err[%s]", err.Error())
    }
    monitor.ServerNum = serverNum
    return monitor
}