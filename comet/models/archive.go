package models

import (
    "github.com/astaxie/beego"
    "fmt"
    "webim/comet/common"
    "github.com/satori/go.uuid"
    "time"
    "encoding/json"
)

type IWorker interface {
    Do()
    Log()
}
type JobWorker struct {
    *Job
    s *Session
}
func NewJobWork(msg Msg, s *Session) *JobWorker {
    reqId := ""
    version := ""
    if v, ok := msg.Data["req_id"]; ok {
        reqId = v.(string)
    }
    if v, ok := msg.Data["version"]; ok {
        version = v.(string)
    }
    job := &Job{
        Version: version,
        ReqID:   reqId,
        TraceID: uuid.NewV4().String(),
        ReqTime: time.Now().Unix(),
        Req:     msg,
    }
    return &JobWorker{job, s}
}
/**
 * 记录日志
 */
func (j *JobWorker) Log(){
    reqJson, _ := json.Marshal(j.Req)
    rspJson, _ := json.Marshal(j.Rsp)
    beego.Info("req:%s, rsp:%s", reqJson, rspJson)
}
func (j *JobWorker) Do() {
    defer j.Log()
    if j.ReqID == "" {
        beego.Error("req_id为空，不做任何处理")
        return
    }
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

func (j *JobWorker) register() {
    if _, ok := j.Req.Data["device_id"]; !ok{

    }
    deviceToken := common.GenerateDeviceToken(j.Req.Data["device_id"].(string), j.Req.Data["appkey"].(string))
    j.Rsp = *NewMsg(TYPE_COMMON_MSG, map[string]interface{}{"code":0, "device_token":deviceToken})
    j.s.Send(j.Job)
}
func (j *JobWorker) leaveRoom() {
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
    j.Rsp = *NewMsg(TYPE_COMMON_MSG, map[string]interface{}{"code":0, "msg":"ok"})
    j.s.Send(j.Job)
}
func (j *JobWorker) joinRoom() {
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
        j.Rsp = *NewMsg(TYPE_COMMON_MSG, data)
        j.s.Send(j.Job)
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
func (j *JobWorker) roomMsg() {
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
        j.Rsp = *NewMsg(TYPE_COMMON_MSG, data)
        j.s.Send(j.Job)
    } else {
        room.Broadcast(&j.Req)
    }
}
func (j *JobWorker) createRoom() {
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
        j.Rsp = *NewMsg(TYPE_COMMON_MSG, data)
        j.s.Send(j.Job)
    }else {
        j.Rsp = j.Req
        j.Rsp.MsgType = TYPE_JOIN_ROOM
        j.s.Send(j.Job)
    }
}
