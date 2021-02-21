package main

import (
	"comet/common"
	"comet/server"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var token string

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		readMsg(c)
	}()
	// 登录加入聊天室
	login(c)

	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			msg := server.Msg{
				Type: server.TYPE_PING,
			}
			msgByte, _ := common.EnJson(msg)
			err := c.WriteMessage(websocket.TextMessage, msgByte)
			if err != nil {
				log.Println("write:", err, "time:", t)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func readMsg(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		// log.Printf("recv: %s", message)
		var tmpMsg server.Msg
		err = common.DeJson(message, &tmpMsg)
		if err != nil {
			log.Printf("json decode err:%+v", err)
			continue
		}
		log.Printf(tmpMsg.Data)
		if tmpMsg.Type == server.TYPE_LOGIN {
			token = tmpMsg.DeviceToken
			joinRoomData := map[string]string{
				"room_id": "1",
			}
			joinRoomDataByte, _ := common.EnJson(joinRoomData)
			joinRoomMsg := server.Msg{
				Type:        server.TYPE_JOIN_ROOM,
				DeviceToken: token,
				Data:        string(joinRoomDataByte),
			}
			joinRoomMsgByte, _ := common.EnJson(joinRoomMsg)
			if err := c.WriteMessage(websocket.TextMessage, joinRoomMsgByte); err != nil {
				log.Printf("join room err:%+v", err)
			}
		} else {

		}

	}
}

func login(c *websocket.Conn) {
	msgData := map[string]string{
		"device_id": common.Uuid(),
		"username":  "demo",
		"password":  "123456",
	}
	msgDataByte, _ := common.EnJson(msgData)
	loginMsg := server.Msg{
		Type:        server.TYPE_LOGIN,
		DeviceToken: "",
		Data:        string(msgDataByte),
	}
	logMsgByte, _ := common.EnJson(loginMsg)

	if err := c.WriteMessage(websocket.TextMessage, logMsgByte); err != nil {
		log.Printf("loging err:%+v", err)
		return
	}
}
