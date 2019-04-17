package server

import (
    "webim/comet/common"
    "github.com/satori/go.uuid"
    "time"
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
    job.Rsp.DeviceToken = s.DeviceToken
    job.Rsp.Version = job.Req.Version
    job.Rsp.ReqId = job.Req.ReqId
    return &JobWorker{job, s}
}
/**
 * 记录日志
 */
func (j *JobWorker) Log(){
    reqJson, _ := common.EnJson(j.Req)
    rspJson, _ := common.EnJson(j.Rsp)
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
    data, err := j.decode(j.Req.Data)
    if err!= nil{
        logs.Error("msg[register decode err] err[%s]", err.Error())
        return
    }
    if _, ok := data["device_id"]; !ok {
        return
    }
    appKey := beego.AppConfig.String("appkey")
    deviceId := data["device_id"].(string)
    deviceToken := common.GenerateDeviceToken(deviceId, appKey)
    j.s.DeviceToken = deviceToken
    j.s.User.DeviceId = deviceId
    if j.s.User.Name == "" {
        j.s.User.Name = deviceId
    }
    //保存token session信息到redis中
    j.s.Server.AddSession(j.s)
    j.Rsp.Type = TYPE_REGISTER
    j.Rsp.DeviceToken = deviceToken
    resData := map[string]interface{}{"code": 0, "device_token": deviceToken}
    resByte, err := common.EnJson(resData)
    if err != nil {
        logs.Error("msg[register encode err] err[%s]", err.Error())
        return
    }
    j.Rsp.Data = string(resByte)
    j.s.Send(&j.Rsp)
}
func (j *JobWorker) leaveRoom() {
    data, err := j.decode(j.Req.Data)
    if err!= nil{
        logs.Error("msg[leaveRoom decode err] err[%s]", err.Error())
        return
    }
    if _, ok := data["room_id"]; !ok {
        logs.Debug("msg[%s]","room_id为空")
        return
    }
    roomId := data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        logs.Error("msg[%s]", err.Error())
        return
    }
    if room != nil {
        room.Leave(j.s)
    }
    j.Rsp.Type = TYPE_ROOM_MSG
    resData := map[string]interface{}{"code":0, "content":"ok"}
    resByte, err := common.EnJson(resData)
    if err != nil {
        logs.Error("msg[leaveRoom encode err] err[%s]", err.Error())
        return
    }
    j.Rsp.Data = string(resByte)
    j.s.Send(&j.Rsp)
}
func (j *JobWorker) joinRoom() {
    data, err := j.decode(j.Req.Data)
    if err!= nil{
        logs.Error("msg[joinRoom decode err] err[%s]", err.Error())
        return
    }
    if _, ok := data["room_id"]; !ok {
        logs.Debug("msg[room_id为空]")
        return
    }
    roomId := data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        logs.Error("msg[%s]", err.Error())
        return
    }
    if room == nil {
        resData := make(map[string]interface{})
        resData["code"] = 1
        resData["content"] = "房间不存在"
        resData["room_id"] = roomId

        resByte, err := common.EnJson(resData)
        if err != nil {
            logs.Error("msg[joinRoom encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Type = TYPE_ROOM_MSG
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    } else {
        res, err := room.Join(j.s)
        if err != nil {
            logs.Error("msg[%s]", err.Error())
            return
        }
        if res {
            resData := make(map[string]interface{})
            resData["code"] = 1
            resData["content"] = j.s.User.Name + "进入房间"
            resByte, err := common.EnJson(resData)
            if err != nil {
                logs.Error("msg[joinRoom encode err] err[%s]", err.Error())
                return
            }
            j.Rsp.Data = string(resByte)
            room.Broadcast(&j.Rsp)
        }
    }
}
func (j *JobWorker) roomMsg() {
    data, err := j.decode(j.Req.Data)
    if err!= nil{
        logs.Error("msg[roomMsg decode err] err[%s]", err.Error())
        return
    }
    if _, ok := data["room_id"]; !ok {
        logs.Warn("msg[room_id为空]")
        return
    }
    roomId := data["room_id"].(string)
    room, err := GetRoom(roomId)
    if err != nil {
        logs.Error("msg[查找房间失败] err[%s]", err.Error())
        return
    }
    if room == nil {
        resData := make(map[string]interface{})
        resData["content"] = "房间不存在"
        resData["room_id"] = roomId

        j.Rsp.Type = TYPE_ROOM_MSG
        resByte, err := common.EnJson(resData)
        if err != nil {
            logs.Error("msg[roomMsg encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    } else {
        j.Rsp = j.Req
        var tmpData map[string]interface{}
        if err = common.DeJson([]byte(j.Req.Data), &tmpData);err !=nil{
            logs.Error("msg[roomMsg DeJson err] err[%s]", err.Error())
            return
        }
        tmpData["uid"] = j.s.User.Id
        tmpData["room_id"] = room.Id
        if TmpRspData, err := common.EnJson(tmpData); err==nil{
            j.Rsp.Data = string(TmpRspData)
            room.Broadcast(&j.Rsp)
        }else{
            logs.Error("msg[roomMsg EnJson err] err[%s]", err)
            return
        }
    }
}
func (j *JobWorker) createRoom() {
    data, err := j.decode(j.Req.Data)
    if err!= nil{
        logs.Error("msg[createRoom decode err] err[%s]", err.Error())
        return
    }
    if _, ok := data["room_id"]; !ok {
        logs.Warn("msg[room_id为空]")
        return
    }
    roomId := data["room_id"].(string)
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
    resData := make(map[string]interface{})
    resData["code"] = 0
    resData["content"] = "创建房间成功"
    j.Rsp.Type = TYPE_CREATE_ROOM
    resByte, err := common.EnJson(resData)
    if err != nil {
        logs.Error("msg[roomMsg encode err] err[%s]", err.Error())
        return
    }
    j.Rsp.Data = string(resByte)
    j.s.Send(&j.Rsp)
}

func (j *JobWorker) decode(jsonStr string) (map[string]interface{}, error) {
    data := map[string]interface{}{}
    err := common.DeJson([]byte(jsonStr), &data)
    if err!=nil{
        return nil, err
    }
    return data, err
}