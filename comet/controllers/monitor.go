package controllers

import (
	"comet/server"
)

type MonitorController struct {
	BaseController
}

func (c MonitorController) Status() {
	c.success(server.Status())
}
