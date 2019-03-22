package controllers

import "webim/comet/models/room"

type RoomController struct {
	BaseController
}

func (c *RoomController) Create() {
	roomId := c.GetString("id")
	room.NewRoom(roomId, "")
	c.success(nil)
}

func (c *RoomController) Delete(){
	roomId := c.GetString("id")
	room.DelRoom(roomId)
	c.success(nil)
}