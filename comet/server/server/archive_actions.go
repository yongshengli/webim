package server

import (
	"comet/common"
	"comet/server/actions"
	"comet/server/base"
	"comet/server/session"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type IWorker interface {
	Do()
	Log()
}
type JobWorker struct {
	*base.Job
	s      *session.Session
	action *actions.Action
}

func NewJobWork(msg base.Msg, s *session.Session) *JobWorker {
	job := &base.Job{
		TraceID: common.LogId(),
		ReqTime: time.Now().Unix(),
		Req:     msg,
	}
	job.Rsp.Encode = "utf-8"
	// job.Rsp.DeviceToken = s.DeviceToken
	job.Rsp.Version = "1.0"
	job.Rsp.ReqId = job.Req.ReqId
	return &JobWorker{job, s, &actions.Action{}}
}

/**
 * 记录日志
 */
func (j *JobWorker) Log() {
	reqJson, _ := common.EnJson(j.Req)
	rspJson, _ := common.EnJson(j.Rsp)
	logs.Info("trace_id[%s] req[%s] rsp[%s] req_time[%d] rsp_time[%d]", j.TraceID, string(reqJson), string(rspJson), j.ReqTime, j.RspTime)
}
func (j *JobWorker) Do() {
	if beego.AppConfig.String("runmode") != "dev" {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("msg[runtime err] err[%v]", r)
			}
		}()
	}
	j.Rsp.Type = j.Req.Type
	defer j.Log()
	switch j.Req.Type {
	case base.TYPE_CREATE_ROOM:
		j.action.CreateRoom(j.s, &j.Req, &j.Rsp)
	case base.TYPE_ROOM_MSG:
		j.action.RoomMsg(j.s, &j.Req, &j.Rsp)
	case base.TYPE_JOIN_ROOM:
		j.action.JoinRoom(j.s, &j.Req, &j.Rsp)
	case base.TYPE_LEAVE_ROOM:
		j.action.LeaveRoom(j.s, &j.Req, &j.Rsp)
	case base.TYPE_REGISTER:
		j.action.Register(j.s, &j.Req, &j.Rsp)
	case base.TYPE_LOGIN:
		j.action.Login(j.s, &j.Req, &j.Rsp)
	case base.TYPE_LOGOUT:
		j.action.Logout(j.s, &j.Req, &j.Rsp)
	case base.TYPE_TRANSPOND:
		j.transpond(j.s, &j.Req, &j.Rsp)
	}
	j.RspTime = time.Now().Unix()
	j.s.Send(j.Rsp)
}

func (j *JobWorker) decode(jsonStr string) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	if err := common.DeJson([]byte(jsonStr), &data); err != nil {
		return nil, err
	}
	return data, nil
}
func (j *JobWorker) transpond() {
	j.TransferReqTime = time.Now().Unix()
	data, err := j.decode(j.Req.Data)
	if err != nil {
		logs.Error("msg[transpond json decode err] err[%s]", err.Error())
		return
	}
	if _, ok := data["url"]; !ok {
		rspData := map[string]interface{}{"code": 1, "content": "url字段不能为空"}
		j.Rsp.Data, _ = common.Map2String(rspData)
		j.s.Send(j.Rsp)
		return
	}
	j.TransferRspTime = time.Now().Unix()
}
