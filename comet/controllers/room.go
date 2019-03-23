package controllers

import "webim/comet/models"

type RoomController struct {
	BaseController
}

func (c *RoomController) Create() {
	roomId := c.GetString("id")
	models.NewRoom(roomId, "")
	c.success(nil)
}

func (c *RoomController) Delete(){
	roomId := c.GetString("id")
	models.DelRoom(roomId)
	c.success(nil)
}