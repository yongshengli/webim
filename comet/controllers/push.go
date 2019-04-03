package controllers

import (
    "webim/comet/models"
)

type PushController struct {
    BaseController
}

func (c *PushController) Unicast(){
    deviceToken := c.GetString("device_token")

    if deviceToken ==""{
        c.error("device_token不能为空")
        return
    }
    msg := models.Msg{MsgType:models.TYPE_NOTICE_MSG, Data:map[string]interface{}{"content":"qwqwq"}}
    models.SessionManager.Unicast(deviceToken, msg)
    c.success(map[string]interface{}{})
}
