package server

import (
    "github.com/astaxie/beego/logs"
    "github.com/astaxie/beego"
    "webim/comet/common"
)

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
func (j *JobWorker) login(){

}

func (j *JobWorker) logout(){

}
