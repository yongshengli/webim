package controllers

import "webim/comet/server"

type RoomController struct {
	BaseController
}

func (c *RoomController) Create() {
	roomId := c.GetString("id")
	server.NewRoom(roomId, "")
	c.success(nil)
}

func (c *RoomController) Delete(){
	roomId := c.GetString("id")
	server.DelRoom(roomId)
	c.success(nil)
}