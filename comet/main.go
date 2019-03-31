// Copyright 2013 Beego Samples authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// This sample is about using long polling and WebSocket to build a web-based chat room based on beego.
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
	go models.RpcServer(":"+rpcPort)
	models.ServerManager.Register(rpcPort)

	beego.Run()
}
