package actions

import (
	"comet/common"
	"comet/models"
	"comet/server/base"
	"comet/server/req2resp"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func (j *Action) Register(s base.Sessioner, req, resp *base.Msg) error {
	reqData := &req2resp.RegisterReq{}
	if err := common.DeJson(req.GetData(), reqData); err != nil {
		logs.Error("msg[register decode err] err[%s]", err.Error())
		return nil
	}
	respData := &req2resp.RegisterResp{}

	if reqData.DeviceId == "" {
		reqData.DeviceId = common.Uuid()
	}
	u := s.GetUser()
	u.DeviceId = respData.DeviceId
	s.SaveLoginState(u)
	respData.DeviceId = reqData.DeviceId

	/**
	同一个设备再次连接时检查此设备是否最近是否登录过
	如果当前连接未登录但是此设备有未过期的session那么将该连接设为已录并关闭历史连接
	**/
	appKey := beego.AppConfig.String("appkey")
	deviceToken := common.GenerateDeviceToken(u.DeviceId, appKey)

	if u.DeviceToken != deviceToken {
		cs := s.GetServer().GetSessionByDeviceToken(deviceToken)
		if cs != nil && cs.DeviceToken != "" && cs != s {
			u.DeviceToken = cs.DeviceToken
			s.SaveLoginState(u)
			cs.Close()
		}
	}
	return nil
}
func (j *Action) login(s base.Sessioner, req, resp *base.Msg) error {
	reqData := &req2resp.LoginReq{}
	if err := common.DeJson(req.GetData(), reqData); err != nil {
		logs.Error("msg[login decode err] err[%s]", err.Error())
		return nil
	}
	respData := &req2resp.LoginResp{}
	if reqData.DeviceId == "" {
		respData.Status = req2resp.Status{
			Code: 1,
			Msg:  "device_id不能为空",
		}
		respByte, _ := common.EnJson(respData)
		resp.SetData(respByte)
		return nil
	}
	if reqData.Username == "" {
		respData.Status = req2resp.Status{
			Code: 1,
			Msg:  "username不能为空",
		}
		respByte, _ := common.EnJson(respData)
		resp.SetData(respByte)
		return nil
	}

	if reqData.Password == "" {
		respData.Status = req2resp.Status{
			Code: 1,
			Msg:  "password不能为空",
		}
		respByte, _ := common.EnJson(respData)
		resp.SetData(respByte)
		return nil
	}
	u, err := models.FindByName(reqData.Username)
	if err != nil {
		respData.Status = req2resp.Status{
			Code: 1,
			Msg:  "用户不存在",
		}
		respByte, _ := common.EnJson(respData)
		resp.SetData(respByte)
		return nil
	}

	if models.CheckPwd(u, reqData.Password) == false {
		respData.Status = req2resp.Status{
			Code: 1,
			Msg:  "用户名或者密码错误",
		}
		respByte, _ := common.EnJson(respData)
		resp.SetData(respByte)
		return nil
	}
	appKey := beego.AppConfig.String("appkey")
	deviceToken := common.GenerateDeviceToken(reqData.DeviceId, appKey)
	user := s.GetUser()
	user.DeviceToken = deviceToken
	user.DeviceId = reqData.DeviceId
	user.Name = reqData.Username
	user.Id = strconv.FormatInt(int64(u.Id), 10)
	s.SaveLoginState(user)

	//保存token session信息到redis中
	// s.Server.AddSession(s)
	// j.Rsp.DeviceToken = deviceToken
	if err != nil {
		logs.Error("msg[register encode err] err[%s]", err.Error())
		return nil
	}
	respData.Status = req2resp.Status{
		Code: 0,
		Msg:  "登录成功",
	}
	respByte, _ := common.EnJson(respData)
	resp.SetData(respByte)
	return nil
}

func (j *Action) logout(req, resp *base.Msg) error {
	return nil
}
