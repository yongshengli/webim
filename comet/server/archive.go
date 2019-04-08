package server

import (
    "webim/comet/common"
    "github.com/satori/go.uuid"
    "time"
    "encoding/json"
    "github.com/astaxie/beego/logs"
    "github.com/astaxie/beego"
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
    job := &Job{
        TraceID: uuid.NewV4().String(),
        ReqTime: time.Now().Unix(),
        Req:     msg,
    }
    job.Rsp.Version = job.Req.Version
    job.Rsp.ReqId = job.Req.ReqId
    return &JobWorker{job, s}
}
/**
 * 记录日志
 */
func (j *JobWorker) Log(){
    reqJson, _ := json.Marshal(j.Req)
    rspJson, _ := json.Marshal(j.Rsp)
    logs.Info("req:%s, rsp:%s", string(reqJson), string(rspJson))
}
func (j *JobWorker) Do() {
    defer j.Log()

    switch j.Req.Type {
    case TYPE_CREATE_ROOM:
        j.createRoom()
    case TYPE_ROOM_MSG:
        j.roomMsg()
    case TYPE_JOIN_ROOM:
        j.joinRoom()
    case TYPE_LEAVE_ROOM:
        j.leaveRoom()
    case TYPE_REGISTER:
        j.register()
    }
}

func (j *JobWorker) register() {
    if _, ok := j.Req.Data["device_id"]; !ok {
        return
    }
    appKey := beego.AppConfig.String("appkey")
    deviceId := j.Req.Data["device_id"].(string)
    deviceToken := common.GenerateDeviceToken(deviceId, appKey)
    j.s.DeviceToken = deviceToken
    j.s.User.DeviceId = deviceId
    if j.s.User.Name == "" {
        j.s.User.Name = deviceId
    }
    //保存token session信息到redis中
    j.s.Server.AddSession(j.s)
    j.Rsp.Type = TYPE_REGISTER
    j.Rsp.Data = map[string]interface{}{"code": 0, "device_token": deviceToken}
    j.s.Send(&j.Rsp)
}
func (j *JobWorker) leaveRoom() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        logs.Debug("msg[%s]","room_id为空")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        logs.Error("msg[%s]", err.Error())
        return
    }
    if room != nil {
        room.Leave(j.s)
    }
    j.Rsp.Type = TYPE_ROOM_MSG
    j.Rsp.Data = map[string]interface{}{"code":0, "content":"ok"}
    j.s.Send(&j.Rsp)
}
func (j *JobWorker) joinRoom() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        logs.Debug("msg[room_id为空]")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        logs.Error("msg[%s]", err.Error())
        return
    }
    if room == nil {
        data := make(map[string]interface{})
        data["content"] = "房间不存在"
        data["room_id"] = roomId
        j.Rsp.Type = TYPE_ROOM_MSG
        j.Rsp.Data = data
        j.s.Send(&j.Rsp)
    } else {
        res, err := room.Join(j.s)
        if err != nil {
            logs.Error("msg[%s]", err.Error())
            return
        }
        if res {
            data := make(map[string]interface{})
            data["content"] = j.s.User.Name + "进入房间"
            j.Rsp.Data = data
            room.Broadcast(&j.Rsp)
        }
    }
}
func (j *JobWorker) roomMsg() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        logs.Warn("msg[room_id为空]")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        logs.Error("msg[查找房间失败] err[%s]", err.Error())
        return
    }
    if room == nil {
        data := make(map[string]interface{})
        data["content"] = "房间不存在"
        data["room_id"] = roomId
        j.Rsp.Type = TYPE_ROOM_MSG
        j.Rsp.Data = data
        j.s.Send(&j.Rsp)
    } else {
        room.Broadcast(&j.Req)
    }
}
func (j *JobWorker) createRoom() {
    if _, ok := j.Req.Data["room_id"]; !ok {
        logs.Warn("msg[room_id为空]")
        return
    }
    roomId := j.Req.Data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        logs.Error("msg[%s]", err.Error())
        return
    }
    if room == nil {
        room, _ = NewRoom(roomId, "")
    }
    _, err = room.Join(j.s)
    if err != nil {
        logs.Error("msg[加入房间失败] err[%s]", err.Error())
        return
    }
    data := make(map[string]interface{})
    data["content"] = "创建房间成功"
    j.Rsp.Type = TYPE_CREATE_ROOM
    j.Rsp.Data = data
    j.s.Send(&j.Rsp)
}
