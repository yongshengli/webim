package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"webim/comet/models"
)
// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	beego.Controller
}
// Join method handles WebSocket requests for WebSocketController.
func (c *WebSocketController) Get() {
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
	models.NewSession(conn, models.SessionManager).Run()
}
func (c *WebIMController) SendMsg(){
	sId := c.GetString("sid")
	msg := models.Msg{MsgType:models.TYPE_COMMON_MSG, Data:map[string]interface{}{"content":"qwqwq"}}
	models.SessionManager.SendMsg(sId, msg)
}