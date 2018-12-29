package room

import (
	"github.com/astaxie/beego"
	"webim/comet/models"
)

type Room struct {
	Id int
	Manager *Manager
	users map[*Session]int
}

func NewRoom(id int, m *Manager) *Room{
	return &Room{Id:id, Manager:m, users:make(map[*Session]int)}
}

func (r *Room) Join(s *Session) bool{
	r.users[s] = 1
	data := make(map[string]interface{})
	data["room_id"] = r.Id
	data["content"] = s.User.Name + "进入房间"
	msg := models.NewMsg(models.TYPE_ROOM_MSG, data)
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
	beego.Debug("room broadcast")
	for session, _ := range r.users{
		session.Send(msg)
	}
	return true
}