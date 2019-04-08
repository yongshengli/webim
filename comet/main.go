package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"webim/comet/controllers"
	"webim/comet/common"
	"webim/comet/server"
)

func main() {
	beego.Info(beego.BConfig.AppName, "start...")

	// Register routers.
	beego.Router("/", &controllers.WebIMController{}, "get:Welcome")
	beego.Router("/webim", &controllers.WebIMController{}, "get:Get")

	beego.Router("/room/create", &controllers.RoomController{}, "post:Create")
	beego.Router("/room/delete", &controllers.RoomController{}, "post:Delete")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{}, "get:Get")
	beego.Router("/push/unicast", &controllers.PushController{}, "post:Unicast")
	beego.Router("/push/broadcast", &controllers.PushController{}, "post:Broadcast")

	common.RedisInit(map[string]string{
		"host": beego.AppConfig.String("redis.host"),
		"port": beego.AppConfig.String("redis.port"),
	})

	// Register template functions.
	beego.AddFuncMap("i18n", i18n.Tr)

	rpcPort := beego.AppConfig.String("rpcport")

	server.Run(common.GetLocalIp(), rpcPort,100, 10000)

	beego.Run()

}
