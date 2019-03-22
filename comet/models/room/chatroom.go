package room

import (
	"github.com/astaxie/beego"
	"webim/comet/models"
	"fmt"
	"webim/comet/common"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"errors"
	"time"
	"log"
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
	//fmt.Println(string(roomJson))
	common.RedisClient.Set(roomKey(id), roomJson, time.Second*3600*24*30)
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
	err = json.Unmarshal(roomJson.([]byte), &room)
	if err!=nil{
		return nil, err
	}
	return &room, nil
}
func DelRoom(id string) (int, error){
	common.RedisClient.Del([]string{roomUserKey(id)})
	return common.RedisClient.Del([]string{roomKey(id)})
}

func roomKey(id string) string{
	return "room:"+id
}

func roomUserKey(roomId string) string{
	return "roomUserList:"+roomId
}

func (r *Room) Users() (map[string]string, error){
	replay, err := common.RedisClient.Do("hgetall", roomUserKey(r.Id))
	if err!=nil{
		return nil, err
	}
	if replay==nil{
		return nil, nil
	}
	return redis.StringMap(replay, err)
}

func (r *Room) Join(ru RUser) (bool, error){
	//r.users[s.Id] = RUser{SId:s.Id, Ip:s.IP, User:*s.User}
	user, err := common.RedisClient.Do("hget", roomUserKey(r.Id), ru.SId)
	if err!= nil {
		log.Println("dddd")
		return false, err
	}
	if user!=nil{
		return true, nil
	}
	jsonStr, err := json.Marshal(ru)
	if err!=nil{
		return false, err
	}
	res, err := redis.Int(common.RedisClient.Do("hset", roomUserKey(r.Id), ru.SId, jsonStr))
	if err!=nil{
		return false, err
	}
	if res<1{
		return false, errors.New("进入房间写入redis失败")
	}
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
	users, err := r.Users()
	if err!=nil{
		return false, err
	}
	for _, user := range users{
		fmt.Println(user)
		//session.Send(msg)
	}
	return true, nil
}