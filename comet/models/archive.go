package models

import (
    "github.com/astaxie/beego"
    "fmt"
    "webim/comet/common"
    "github.com/satori/go.uuid"
    "time"
)

func NewJob(msg Msg) *Job {
    return &Job{Version: msg.Data["version"].(string),
        ReqID: msg.Data["req_id"].(string),
        TraceID: uuid.NewV4().String(),
        ReqTime: time.Now().Unix(),
        Req: msg,
        s: nil}
}
func (j *Job) Run(s *Session) {
    j.s = s
    switch j.Req.MsgType {
    case TYPE_CREATE_ROOM:
        j.createRoom()
    case TYPE_ROOM_MSG:
        j.roomMsg()
    case TYPE_JOIN_ROOM:
        j.joinRoom()
    case TYPE_LEAVE_ROOM:
        j.leaveRoom()
    }
}
func (j *Job) register() {
    if _, ok := j.Req.Data["device_id"]; !ok{

    }
    deviceToken := common.GenerateDeviceToken(j.Data["device_id"].(string), j.Data["appkey"].(string))
    j.s.Send(NewMsg(TYPE_COMMON_MSG, map[string]interface{}{"device_token":deviceToken}))
}
func (j *Job) leaveRoom() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        fmt.Println("room_id 为空")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room != nil {
        room.Leave(j.s.Id)
    }
}
func (j *Job) joinRoom() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        beego.Warn("room_id 为空")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room == nil {
        data := make(map[string]interface{})
        data["content"] = "房间不存在"
        j.s.Send(NewMsg(TYPE_COMMON_MSG, data))
    } else {
        ru := RUser{SId: j.s.Id, User: *j.s.User, Addr: j.s.Addr}
        res, err := room.Join(RUser{SId: j.s.Id, User: *j.s.User, Addr: j.s.Addr})
        j.s.RoomId = roomId
        if err != nil {
            beego.Error(err)
        }
        if res {
            data := make(map[string]interface{})
            data["room_id"] = room.Id
            data["content"] = ru.User.Name + "进入房间"
            msg := NewMsg(TYPE_ROOM_MSG, data)
            room.Broadcast(msg)
        }
    }
}
func (j *Job) roomMsg() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        beego.Warn("room_id 为空")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room == nil {
        data := make(map[string]interface{})
        data["content"] = "房间不存在"
        j.s.Send(NewMsg(TYPE_COMMON_MSG, data))
    } else {
        room.Broadcast(&j.Req)
    }
}
func (j *Job) createRoom() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        beego.Warn("room_id 为空")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        beego.Error(err)
        return
    }
    if room == nil {
        NewRoom(roomId, "")
        data := make(map[string]interface{})
        data["content"] = "创建房间成功"
        j.s.Send(NewMsg(TYPE_COMMON_MSG, data))
    }
    j.Rsp = j.Req
    j.Rsp.MsgType = TYPE_JOIN_ROOM
    j.s.Send(&j.Rsp)
}
