package room

import (
	"github.com/astaxie/beego"
	"webim/comet/models"
	"fmt"
	_ "webim/comet/common"
	"webim/comet/common"
	"encoding/json"
	"comet/src/rrx/redigo/redis"
	"errors"
)

type RUser struct {
	SId  string `json:"sid"`
	IP   string `json:"ip"`   //sid 所在机器ip
	User User   `json:"user"` //用户数据
}
type Room struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
}

func NewRoom(id string, name string) (*Room, error){
	room := &Room{Id:id, Name:name}
	roomJson, err := json.Marshal(room)
	if err!=nil{
		beego.Error(err)
		return nil, err
	}
	common.RedisClient.Set(roomKey(id), string(roomJson), 0)
	return room, nil
}

func RoomList() []Room{
	return []Room{}
}

func GetRoom(id string) (*Room, error){
	roomJson, err := common.RedisClient.Get(roomKey(id))
	if err!=nil{
		return nil, err
	}
	var room Room
	err = json.Unmarshal([]byte(roomJson.(string)), &room)
	if err!=nil{
		return nil, err
	}
	return &room, nil
}
func DelRoom(id string) (int, error){
	return common.RedisClient.Del([]string{roomKey(id)}), nil
}

func roomKey(id string) string{
	return "room:"+id
}

func roomUserKey(roomId string) string{
	return "roomUserList:"+roomId
}

func (r *Room) Users() map[string]RUser{
	common.RedisClient.Do("hgetall", roomUserKey(r.Id))
	return map[string]RUser{}
}

func (r *Room) Join(ru RUser) (bool, error){
	//r.users[s.Id] = RUser{SId:s.Id, Ip:s.IP, User:*s.User}
	tMap, err := redis.StringMap(common.RedisClient.Do("hget", roomUserKey(r.Id), ru.SId))
	if err!= nil {
		return false, err
	}
	if _, ok := tMap["Sid"]; ok{
		return true, nil
	}
	jsonStr, err := json.Marshal(ru)
	if err!=nil{
		return false, err
	}
	res, err := redis.Int(common.RedisClient.Do("hset", roomUserKey(r.Id), ru.SId, string(jsonStr)))
	if err!=nil{
		return false, err
	}
	if res<1{
		return false, errors.New("进入房间写入redis失败")
	}
	data := make(map[string]interface{})
	data["room_id"] = r.Id
	data["content"] = ru.User.Name + "进入房间"
	msg := models.NewMsg(models.TYPE_ROOM_MSG, data)
	r.Broadcast(msg)
	return true, nil
}

func (r *Room) Leave(ru RUser) (bool, error){
	_, err := common.RedisClient.Do("hdel", roomUserKey(r.Id), ru.SId)
	if err!=nil{
		return false, err
	}
	userNum, err := redis.Int(common.RedisClient.Do("hlen", roomUserKey(r.Id)))
	if err!= nil{
		return false, err
	}
	if userNum<1 {
		DelRoom(r.Id)
		beego.Debug("房间%d内用户为空删除房间", r.Id)
	}
	return true, nil
}
// This function handles all incoming chan messages.

//房间内广播
func (r *Room) Broadcast(msg *models.Msg) (bool, error){
	msg.Data["room_id"] = r.Id
	beego.Debug("room broadcast")
	users := r.Users()
	for _, user := range users{
		fmt.Println(user)
		//session.Send(msg)
	}
	return true, nil
}