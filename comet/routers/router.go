package routers

import (
	"webim/comet/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// Register routers.
	beego.Router("/", &controllers.WebIMController{}, "get:Welcome")
	// Indicate AppController.Join method to handle POST requests.
	beego.Router("/room/create", &controllers.WebIMController{}, "post:Join")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})
}
