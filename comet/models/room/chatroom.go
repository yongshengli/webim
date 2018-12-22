package room

import (
	"webim/comet/models"
)

type Room struct {
	Id int
	Manager *Manager
	users map[*Session]int
}

type JoinRoomData struct {
	RoomId int
	User User
}

func (r *Room) Join(s *Session) bool{
	r.users[s] = 1
	data := make(map[string]interface{})
	data["room_id"] = r.Id
	data["content"] = s.User.Name + "进入房间"
	msg := models.New(models.TYPE_ROOM_MSG, data)
	r.Broadcast(msg)
	return true
}

func (r *Room) Leave(s *Session) bool{
	if _, ok := r.users[s]; !ok{
		return true
	}
	delete(r.users, s)
	return true
}
// This function handles all incoming chan messages.

//房间内广播
func (r *Room) Broadcast(msg *models.Msg) bool{
	msg.Data["room_id"] = r.Id

	for session, _ := range r.users{
		session.Send(msg)
	}
	return true
}