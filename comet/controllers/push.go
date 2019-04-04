package controllers

import (
    "webim/comet/models"
)

type PushController struct {
    BaseController
}

func (c *PushController) Unicast(){
    params := c.Data["params"].(map[string]interface{})

    if deviceToken, ok := params["device_token"]; !ok || deviceToken==""{
        c.error("device_token不能为空")
        return
    }
    if _, tOk := params["msg"].(map[string]interface{}); !tOk{
        c.error("msg为空或者msg格式错误")
        return
    }
    deviceToken := params["device_token"].(string)
    msg := models.Msg{MsgType:models.TYPE_NOTICE_MSG, Data:map[string]interface{}{"content":"qwqwq"}}
    _, err := models.SessionManager.Unicast(deviceToken, msg)
    if err!=nil{
        c.error(err.Error())
        return
    }
    c.success(nil)
}

func (c *PushController) Broadcast(){

}
