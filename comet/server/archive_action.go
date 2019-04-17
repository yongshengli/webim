package server

import (
    "webim/comet/common"
    "github.com/satori/go.uuid"
    "time"
    "github.com/astaxie/beego/logs"
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
    logs.Info("trace_id[%s] req[%s] rsp[%s] req_time[%s] rsp_time[%s]", j.TraceID, string(reqJson), string(rspJson), j.ReqTime, j.RspTime)
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
    err := common.DeJson([]byte(jsonStr), &data)
    if err!=nil{
        return nil, err
    }
    return data, err
}
func (j *JobWorker) transpond(){

}