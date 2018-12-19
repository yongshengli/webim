package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"webim/room/models"
)

type User struct {
	Uuid string
	Id int
	Name string
	Info map[string]interface{}
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

type JoinRoomData struct {
	RoomId int
	User User
}
var (
	// Channel for new join users.
	joinRoomChan = make(chan JoinRoomData, 100)
	// Channel for exit users.
	leaveRoomChan = make(chan JoinRoomData, 100)
	// Send events here to publish them.
	reqChan = make(chan models.Msg, 100)
	repChan = make(chan models.Msg, 100)
	rooms = make(map[int]map[*websocket.Conn]int)
	userMap = make(map[int]*websocket.Conn)
	connMap = make(map[*websocket.Conn]User)
)

func saveConn(conn *websocket.Conn){
	connMap[conn] = User{Conn:conn}
}

func disConn(conn *websocket.Conn){
	conn.Close()
	if u, ok := connMap[conn]; ok{
		if _, ok2 := userMap[u.Id]; ok2{
			delete(userMap, u.Id)
		}
		delete(connMap, conn)
	}
}

func joinRoom(roomId int, u User) bool{
	connMap[u.Conn] = u
	if _, ok :=rooms[roomId]; !ok{
		rooms[roomId] = make(map[*websocket.Conn]int)
	}
	rooms[roomId][u.Conn] = u.Id
	return true
}

func leaveRoom(roomId int, u User) bool{
	if _, ok := rooms[roomId][u.Conn]; !ok{
		return true
	}
	delete(rooms[roomId], u.Conn)
	return true
}
// This function handles all incoming chan messages.
func chatroom() {
	for {
		select {
		case sub := <- joinRoomChan:
			if _, ok := rooms[sub.RoomId][sub.User.Conn]; ok{
				beego.Info("Old user:", sub.User.Name, ";WebSocket:", sub.User.Conn != nil)
			} else {
				joinRoom(sub.RoomId, sub.User) // Add user to the end of list.
				// Publish a JOIN event.
				//publish <- newEvent(models.EVENT_JOIN,  EventData{})
				beego.Info("New user:", sub.User.Name, ";WebSocket:", sub.User.Conn != nil)
			}
		case req := <- reqChan:
			// Notify waiting list.
			if req.MsgType==models.MSG_ROOM{
				repChan <- req
			}
		case rep := <- repChan:
			if rep.MsgType==models.MSG_ROOM{
				roomBroadcast(rep)
			}
		case unsub := <-leaveRoomChan:
			leaveRoom(unsub.RoomId, unsub.User)
		}
	}
}

func init() {
	go chatroom()
}

func isUserExist(uid int) bool {
	if _, ok := userMap[uid]; ok{
		return true;
	}
	return false
}

//unicast
func unicast(uid int, msg models.Msg) bool{
	if conn, ok := userMap[uid]; ok {
		data, err := json.Marshal(msg)
		if err != nil {
			beego.Error("Fail to marshal event:", err)
			return false
		}
        if conn.WriteMessage(websocket.TextMessage, data) !=nil {
			disConn(conn)
		}
	}
	return false
}

//broadcasts messages to WebSocket users.
func broadcast(msg models.Msg) bool{
	data, err := json.Marshal(msg)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return false
	}
	for _, u := range connMap{
		// Immediately send event to WebSocket users.
		if u.Conn != nil {
			if u.Conn.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				disConn(u.Conn)
			}
		}
	}
	return true
}

//房间内广播
func roomBroadcast(roomId int, msg models.Msg) bool{
	if _, ok := rooms[roomId]; !ok{
		return false
	}
	data, err := json.Marshal(msg)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return false
	}
	for conn, _ := range rooms[roomId]{
		// Immediately send event to WebSocket users.
		if conn != nil {
			if conn.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected. delete from room
				leaveRoom(roomId, connMap[conn])
				disConn(conn)
			}
		}
	}
	return true
}