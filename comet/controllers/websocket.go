package controllers

import (
	"webim/comet/models/room"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	BaseController
}

// Get method handles GET requests for WebSocketController.
func (c *WebSocketController) Get() {
	// Safe check.
	uname := c.GetString("uname")
	if len(uname) == 0 {
		c.Redirect("/", 302)
		return
	}

	c.TplName = "websocket.html"
	c.Data["IsWebSocket"] = true
	c.Data["UserName"] = uname
	c.Render()
}

// Join method handles WebSocket requests for WebSocketController.
func (c *WebSocketController) Join() {
	ws := websocket.Upgrader{}
	// Upgrade from http request to WebSocket.
	conn, err := ws.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	room.NewSession(conn, room.SessionManager).Run()
}