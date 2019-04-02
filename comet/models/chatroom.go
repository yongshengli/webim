package models

import (
	"github.com/astaxie/beego"
	"webim/comet/common"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"errors"
	"time"
)

type RUser struct {
	DeviceToken string `json:"device_token"`
	IP          string `json:"ip"`   //sid 所在机器ip
	User        User   `json:"user"` //用户数据
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
	if roomJson==nil{
		return nil, nil
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
	return "comet:room:"+id
}

func roomUserKey(roomId string) string{
	return "comet:roomUserList:"+roomId
}

func (r *Room) Users() (map[string]interface{}, error){
	replay, err := common.RedisClient.Do("hgetall", roomUserKey(r.Id))
	if err != nil {
		return nil, err
	}
	if replay == nil {
		return nil, nil
	}
	tmap, err := redis.StringMap(replay, err)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{}
	for dt, st := range tmap {
		ru := RUser{}
		json.Unmarshal([]byte(st), &ru)
		res[dt] = ru
	}
	return res, nil
}

func (r *Room) Join(s *Session) (bool, error){
	//r.users[s.Id] = RUser{SId:s.Id, Ip:s.IP, User:*s.User}
	//退出旧房间
	if s.RoomId != "" {
		beego.Debug("用户%s退出旧房间%s", s.DeviceToken, s.RoomId)
		oldR, _ := GetRoom(s.RoomId)
		if oldR != nil {
			oldR.Leave(s)
		}
	}
	s.RoomId = r.Id
	ru := RUser{DeviceToken: s.DeviceToken, User: *s.User, IP: s.IP}
	user, err := common.RedisClient.Do("hget", roomUserKey(r.Id), ru.DeviceToken)
	if err!= nil {
		return false, err
	}
	if user!=nil{
		return true, nil
	}
	jsonStr, err := json.Marshal(ru)
	if err!=nil{
		return false, err
	}
	res, err := redis.Int(common.RedisClient.Do("hset", roomUserKey(r.Id), ru.DeviceToken, jsonStr))
	if err!=nil{
		return false, err
	}
	if res<1{
		return false, errors.New("进入房间写入redis失败")
	}
	return true, nil
}

func (r *Room) Leave(s *Session) (bool, error){
	if len(s.DeviceToken)<1{
		return false, errors.New("session.DeviceToken为空")
	}
	_, err := common.RedisClient.Do("hdel", roomUserKey(r.Id), s.DeviceToken)
	if err!=nil{
		return false, err
	}
	s.RoomId = ""
	userNum, err := redis.Int(common.RedisClient.Do("hlen", roomUserKey(r.Id)))
	if err!= nil{
		return false, err
	}

	if userNum < 1 {
		DelRoom(r.Id)
		beego.Debug("房间%d内用户为空删除房间", r.Id)
	}
	return true, nil
}
// This function handles all incoming chan messages.

//房间内广播
func (r *Room) Broadcast(msg *Msg) (bool, error){
	msg.Data["room_id"] = r.Id
	beego.Debug("msg[room broadcast]")
	users, err := r.Users()
	if err!=nil{
		return false, err
	}
	//beego.Debug("msg[room中的用户] users[%s]", users)
	for _, user := range users{
		//session.Send(msg)
        SessionManager.Unicast(user.(RUser).DeviceToken, *msg)
	}
	return true, nil
}