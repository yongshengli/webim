package controllers

import (
    "webim/comet/models"
)

type PushController struct {
    BaseController
}

func (c *PushController) Unicast() {
    params := c.Data["params"].(map[string]interface{})

    if deviceToken, ok := params["device_token"]; !ok || deviceToken == "" {
        c.error("device_token不能为空")
        return
    }
    if _, tOk := params["msg"].(map[string]interface{}); !tOk {
        c.error("msg为空或者msg格式错误")
        return
    }
    deviceToken := params["device_token"].(string)
    msg := models.Map2Msg(params["msg"].(map[string]interface{}))

    _, err := models.SessionManager.Unicast(deviceToken, msg)
    if err != nil {
        c.error(err.Error())
        return
    }
    c.success(nil)
}

func (c *PushController) Broadcast() {
    params := c.Data["params"].(map[string]interface{})
    if _, ok := params["msg"].(map[string]interface{}); !ok {
        c.error("msg为空或者msg格式错误")
        return
    }
    msg := models.Map2Msg(params["msg"].(map[string]interface{}))

    _, err := models.SessionManager.Broadcast(msg)
    if err != nil {
        c.error(err.Error())
        return
    }
    c.success(nil)
}
