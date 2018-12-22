package controllers

import "webim/comet/models/room"

type RoomController struct {
	BaseController
}

func (c *RoomController) Create() {
	roomId, err := c.GetInt("id")
	if err !=nil{
		c.error(err.Error())
		return
	}
	room.SessionManager.AddRoom(roomId)
	c.success(nil)
}

func (c *RoomController) Delete(){
	roomId, err := c.GetInt("id")
	if err !=nil{
		c.error(err.Error())
		return
	}
	room.SessionManager.DelRoom(roomId)
	c.success(nil)
}