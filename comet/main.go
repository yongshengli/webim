package main

import (
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"webim/comet/controllers"
	"webim/comet/models"
	"webim/comet/common"
)

func main() {
	beego.Info(beego.BConfig.AppName, "start...")

	// Register routers.
	beego.Router("/", &controllers.WebIMController{}, "get:Welcome")
	beego.Router("/webim", &controllers.WebIMController{})

	beego.Router("/room/create", &controllers.RoomController{}, "post:Create")
	beego.Router("/room/delete", &controllers.RoomController{}, "post:Delete")

	// WebSocket.
	beego.Router("/ws", &controllers.WebSocketController{})

	common.RedisInit(map[string]string{
		"host": beego.AppConfig.String("redis.host"),
		"port": beego.AppConfig.String("redis.port"),
	})

	// Register template functions.
	beego.AddFuncMap("i18n", i18n.Tr)
	rpcPort := beego.AppConfig.String("rpcport")
	go models.RpcServer(rpcPort)
	models.ServerManager.Register(rpcPort)

	beego.Run()
}
