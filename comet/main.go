package main

import (
	"comet/common"
	"comet/controllers"
	"comet/models"
	"comet/service/server"
	"math/bits"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

func main() {

	bits.OnesCount()

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
	beego.Router("/monitor/status", &controllers.MonitorController{}, "get:Status")

	beego.Router("/user/register", &controllers.UserController{}, "get,post:Register")

	common.RedisInit(map[string]string{
		"host": beego.AppConfig.String("redis.host"),
		"port": beego.AppConfig.String("redis.port"),
	})
	err := models.ConnectMysql(
		beego.AppConfig.String("mysql.host"),
		beego.AppConfig.String("mysql.port"),
		beego.AppConfig.String("mysql.user"),
		beego.AppConfig.String("mysql.pass"),
		beego.AppConfig.String("mysql.db"),
	)
	if err != nil {
		panic(err)
	}
	// Register template functions.
	beego.AddFuncMap("i18n", i18n.Tr)

	rpcPort := beego.AppConfig.String("rpcport")

	server.Run(common.GetLocalIp(), rpcPort, 100, 10000)

	beego.Run()

}
