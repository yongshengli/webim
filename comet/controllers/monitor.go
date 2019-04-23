package controllers

import (
	"webim/comet/server"
)

type MonitorController struct{
	BaseController
}

func (c MonitorController) Status(){
	c.success(server.Status())
}