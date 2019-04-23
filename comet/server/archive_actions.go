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
    job.Rsp.Encode = "utf-8"
    job.Rsp.DeviceToken = s.DeviceToken
    job.Rsp.Version = "1.0"
    job.Rsp.ReqId = job.Req.ReqId
    return &JobWorker{job, s}
}
/**
 * 记录日志
 */
func (j *JobWorker) Log(){
    reqJson, _ := common.EnJson(j.Req)
    rspJson, _ := common.EnJson(j.Rsp)
    logs.Info("trace_id[%s] req[%s] rsp[%s] req_time[%d] rsp_time[%d]", j.TraceID, string(reqJson), string(rspJson), j.ReqTime, j.RspTime)
}
func (j *JobWorker) Do() {
    if beego.AppConfig.String("runmode")!="dev" {
        defer func() {
            if r := recover(); r != nil {
                logs.Error("msg[runtime err] err[%v]", r)
            }
        }()
    }
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
    case TYPE_LOGIN:
        j.login()
    case TYPE_LOGOUT:
        j.logout()
    case TYPE_TRANSPOND:
        j.transpond()
    }
    j.RspTime = time.Now().Unix()
}

func (j *JobWorker) decode(jsonStr string) (map[string]interface{}, error) {
    data := map[string]interface{}{}
    if err := common.DeJson([]byte(jsonStr), &data);err!=nil{
        return nil, err
    }
    return data, nil
}
func (j *JobWorker) transpond(){
    j.TransferReqTime = time.Now().Unix()
    data, err := j.decode(j.Req.Data)
    if err !=nil {
        logs.Error("msg[transpond json decode err] err[%s]", err.Error())
        return
    }
    if _, ok := data["url"]; !ok{
        rspData := map[string]interface{}{"code":1, "content":"url字段不能为空"}
        byteData, err := common.EnJson(rspData)
        if err !=nil {
            return
        }
        j.Rsp.Data = string(byteData)
        j.s.Send(&j.Rsp)
        return
    }
    j.TransferRspTime = time.Now().Unix()
}