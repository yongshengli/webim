package server

import (
	"webim/comet/common"
	"webim/comet/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func (j *JobWorker) register() {
	data, err := j.decode(j.Req.Data)
	if err != nil {
		logs.Error("msg[register decode err] err[%s]", err.Error())
		return
	}
	if _, ok := data["device_id"]; !ok {
		return
	}

	j.Rsp.Type = TYPE_REGISTER
	j.s.User.DeviceId = data["device_id"].(string)
	j.Rsp.Data, _ = common.Map2String(map[string]interface{}{"code": 0})
	j.s.Send(j.Rsp)
}
func (j *JobWorker) login() {
	data, err := j.decode(j.Req.Data)
	if err != nil {
		logs.Error("msg[register decode err] err[%s]", err.Error())
		return
	}

	j.Rsp.Type = TYPE_LOGIN
	if deviceId, ok := data["device_id"]; !ok || deviceId == "" {
		j.Rsp.Data, _ = common.Map2String(map[string]interface{}{"code": 1, "content": "device_id不能为空"})
		j.s.Send(j.Rsp)
		return
	}
	if userName, ok := data["username"]; !ok || userName == "" {
		j.Rsp.Data, _ = common.Map2String(map[string]interface{}{"code": 1, "content": "username不能为空"})
		j.s.Send(j.Rsp)
		return
	}

	if pass, ok := data["password"]; !ok || pass == "" {
		j.Rsp.Data, _ = common.Map2String(map[string]interface{}{"code": 1, "content": "password不能为空"})
		j.s.Send(j.Rsp)
		return
	}
	deviceId := data["device_id"].(string)
	userName := data["username"].(string)
	pass := data["password"].(string)
	u, err := models.FindByName(userName)
	if err != nil {
		j.Rsp.Data, _ = common.Map2String(map[string]interface{}{"code": 1, "content": "用户不存在"})
		j.s.Send(j.Rsp)
		return
	}
	if models.CheckPwd(u, pass) == false {
		j.Rsp.Data, _ = common.Map2String(map[string]interface{}{"code": 1, "content": "用户名或者密码错误"})
		j.s.Send(j.Rsp)
		return
	}
	appKey := beego.AppConfig.String("appkey")
	deviceToken := common.GenerateDeviceToken(deviceId, appKey)
	j.s.DeviceToken = deviceToken
	j.s.User.DeviceId = deviceId
	if j.s.User.Name == "" {
		j.s.User.Name = userName
	}
	//保存token session信息到redis中
	j.s.Server.AddSession(j.s)
	j.Rsp.DeviceToken = deviceToken
	if err != nil {
		logs.Error("msg[register encode err] err[%s]", err.Error())
		return
	}
	j.Rsp.Data, _ = common.Map2String(map[string]interface{}{"code": 0, "content": "登录成功"})
	j.s.Send(j.Rsp)
}

func (j *JobWorker) logout() {

}
