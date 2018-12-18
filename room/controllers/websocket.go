package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"webim/room/models"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	BaseController
}

// Get method handles GET requests for WebSocketController.
func (this *WebSocketController) Get() {
	// Safe check.
	uname := this.GetString("uname")
	if len(uname) == 0 {
		this.Redirect("/", 302)
		return
	}

	this.TplName = "websocket.html"
	this.Data["IsWebSocket"] = true
	this.Data["UserName"] = uname
}

// Join method handles WebSocket requests for WebSocketController.
func (this *WebSocketController) Join() {
	upgrader := websocket.Upgrader{}
	// Upgrade from http request to WebSocket.
	ws, err := upgrader.Upgrade(this.Ctx.ResponseWriter, this.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	conn(ws)
	defer disConn(ws)

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- models.NewEvent(models.EVENT_MESSAGE, string(p))
	}
}