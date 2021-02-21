package room

import (
	"comet/common"
	"comet/models"
	"comet/server/base"
	"errors"
	"reflect"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

//Room 房间结构定义
type Room struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//NewRoom 新建房间
func NewRoom(id string, name string) (*Room, error) {
	room := &Room{Id: id, Name: name}
	roomJson, err := common.EnJson(room)
	if err != nil {
		logs.Error("msg[room struct json encode err] err[%s]", err.Error())
		return nil, err
	}
	commands := make([]common.RedisCommands, 2)
	commands[0] = common.RedisCommands{CommandName: "SETEX", Args: []interface{}{roomKey(id), int(base.ROOM_LIVE_TIME / time.Second), roomJson}}
	commands[1] = common.RedisCommands{CommandName: "ZADD", Args: []interface{}{roomZsetKey(), time.Now().Unix(), id}}
	res := common.RedisClient.Pipeline(commands)
	//fmt.Println(res)
	for i := 0; i < len(commands); i++ {
		if res[i]["err"] != nil {
			logs.Error("msg[新建Room失败] err[%s] result[%v]", res[i]["err"].(error), res[i]["reply"])
			return nil, err
		}
	}
	return room, nil
}

//RoomList 房间列表
func RoomList(start, stop int) (map[string]string, error) {
	res, err := common.RedisClient.Do("ZREVRANGE", roomZsetKey(), start, stop, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	if reflect.ValueOf(res).Kind() == reflect.Slice {
		result, err := redis.StringMap(res, err)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("redis 返回的数据结构错误")
}

//TotalRoom 统计总房间数
func TotalRoom() (num int, err error) {
	return redis.Int(common.RedisClient.Do("ZCARD", roomZsetKey()))
}

//GetRoom 获取房间
func GetRoom(id string) (*Room, error) {
	roomRedisKey := roomKey(id)
	commands := make([]common.RedisCommands, 2)
	commands[0] = common.RedisCommands{CommandName: "GET", Args: []interface{}{roomRedisKey}}
	commands[1] = common.RedisCommands{CommandName: "TTL", Args: []interface{}{roomRedisKey}}
	res := common.RedisClient.Pipeline(commands)
	if res[0]["err"] != nil {
		return nil, res[0]["err"].(error)
	}
	if res[1]["err"] != nil {
		return nil, res[1]["err"].(error)
	}
	//fmt.Println(res)
	var room Room
	if res[0]["reply"] == nil {
		return nil, nil
	}
	if err := common.DeJson(res[0]["reply"].([]byte), &room); err != nil {
		logs.Error("msg[GetRoom解析room redis data err] err[%s]", err.Error())
		return nil, err
	}

	if ttl := res[1]["reply"].(int64); ttl < 3600*3 {
		_, err := common.RedisClient.Expire(roomRedisKey, base.ROOM_LIVE_TIME)
		if err != nil {
			logs.Error("msg[延长room data redis 时间失败] err[%s]", err.Error())
		}
	}
	return &room, nil
}

//DelRoom 删除房间
func DelRoom(id string) (int, error) {
	commands := make([]common.RedisCommands, 3)
	commands[0] = common.RedisCommands{CommandName: "DEL", Args: []interface{}{roomUserKey(id)}}
	commands[1] = common.RedisCommands{CommandName: "DEL", Args: []interface{}{roomKey(id)}}
	commands[2] = common.RedisCommands{CommandName: "ZREM", Args: []interface{}{roomZsetKey(), id}}
	res := common.RedisClient.Pipeline(commands)
	for i := 0; i < len(commands); i++ {
		if res[i]["err"] != nil {
			return 0, res[i]["err"].(error)
		}
	}
	return redis.Int(res[1]["reply"], nil)
}

func roomKey(id string) string {
	return "comet:room:" + id
}
func roomZsetKey() string {
	return "comet:roomZset"
}
func roomUserKey(roomId string) string {
	return "comet:roomUserList:" + roomId
}

//Users 房间内全部用户的token
func (r *Room) Users(start, end int) ([]string, error) {
	tokenSlice, err := redis.Strings(common.RedisClient.Do("zrange", roomUserKey(r.Id), start, end))
	if err != nil {
		return nil, err
	}
	return tokenSlice, nil
}

//Join 加入房间
func (r *Room) Join(s base.Sessioner) (bool, error) {
	//r.users[s.Id] = RUser{SId:s.Id, Ip:s.IP, User:*s.User}
	//退出旧房间
	if s.RoomId != "" && s.RoomId != r.Id {
		beego.Debug("msg[用户%s退出旧房间%s]", s.DeviceToken, s.RoomId)
		oldR, _ := GetRoom(s.RoomId)
		if oldR != nil {
			oldR.Leave(s)
		}
	}
	s.RoomId = r.Id
	// ru := RUser{DeviceToken: , User: *s.User}
	res, err := redis.Int(common.RedisClient.Do("zadd", roomUserKey(r.Id), time.Now().Unix(), s.DeviceToken))
	if err != nil {
		return false, err
	}
	if res < 1 {
		return false, errors.New("进入房间写入redis失败")
	}
	return true, nil
}

//Leave 离开房间
func (r *Room) Leave(s base.Sessioner) (bool, error) {
	if len(s.DeviceToken) < 1 {
		return false, errors.New("session.DeviceToken为空")
	}

	if _, err := common.RedisClient.Do("zrem", roomUserKey(r.Id), s.DeviceToken); err != nil {
		return false, err
	}
	s.RoomId = ""
	userNum, err := redis.Int(common.RedisClient.Do("zcard", roomUserKey(r.Id)))
	if err != nil {
		return false, err
	}

	if userNum < 1 {
		DelRoom(r.Id)
		beego.Debug("msg[房间%d内用户为空删除房间]", r.Id)
	}
	return true, nil
}

//Broadcast 房间内广播
func (r *Room) Broadcast(msgData *models.RoomMsg) (bool, error) {
	msgStr, _ := common.EnJson(msgData)
	msg := base.Msg{
		Type: base.TYPE_ROOM_MSG,
		Data: string(msgStr),
	}

	if roomMsgId, err := SaveRoomMsg(r.Id, msgData); err != nil || roomMsgId == 0 {
		logs.Error("msg[room msg 写入数据库失败]")
	}
	jsonByte, err := common.EnJson(msgData)
	if err != nil {
		return false, err
	}
	msg.Data = string(jsonByte)
	beego.Debug("msg[room broadcast]")
	start := 0
	pageSize := 1000
	for {
		end := start + pageSize - 1
		users, err := r.Users(start, start+pageSize)
		start = end + 1
		if err != nil {
			return false, err
		}
		if len(users) < 1 {
			break
		}
		//beego.Debug("msg[room中的用户] users[%s]", users)
		for _, user := range users {
			//session.Send(msg)
			tmsg := msg
			IMServer.Unicast(user, tmsg)
		}
	}
	return true, nil
}

//SaveRoomMsg 保存聊天记录到
func SaveRoomMsg(roomId string, roomMsg *models.RoomMsg) (uint64, error) {
	_, err := models.SaveMsg2Redis(roomMsg)
	go func() {
		for i := 0; i < 3; i++ {
			res := models.InsertRoomMsg(*roomMsg)
			if res != nil && len(res.GetErrors()) == 0 {
				return
			}
			if i == 2 {
				logs.Error("msg[i:%d,写入mysql失败err:%+v]", i, res.GetErrors())
			}
		}
	}()
	return roomMsg.MsgId, err
}
