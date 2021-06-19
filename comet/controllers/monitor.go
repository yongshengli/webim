package controllers

import (
	"comet/service/server"
)

type MonitorController struct {
	BaseController
}

func (c MonitorController) Status() {
	c.success(server.Status())
}
