package server

import (
    "github.com/astaxie/beego"
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
    "errors"
    "time"
    "github.com/astaxie/beego/logs"
    "reflect"
    "webim/comet/models"
    "strconv"
)

type RUser struct {
    DeviceToken string `json:"device_token"`
    IP          string `json:"ip"`   //sid 所在机器ip
    User        User   `json:"user"` //用户数据
}
type Room struct {
    Id   string `json:"id"`
    Name string `json:"name"`
}

func NewRoom(id string, name string) (*Room, error) {
    room := &Room{Id: id, Name: name}
    roomJson, err := common.EnJson(room)
    if err != nil {
        logs.Error("msg[room struct json encode err] err[%s]", err.Error())
        return nil, err
    }
    res, err := common.RedisClient.Multi(func(conn redis.Conn) {
        conn.Send("SETEX", roomKey(id), int(ROOM_LIVE_TIME/time.Second), roomJson)
        conn.Send("ZADD", roomZsetKey(), time.Now().Unix(), id)
    })
    //fmt.Println(res)
    if err != nil {
        logs.Error("msg[新建Room失败] err[%s] result[%s]", err.Error(), res)
        return nil, err
    }
    return room, nil
}

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

func TotalRoom() (num int, err error) {
    return redis.Int(common.RedisClient.Do("ZCARD", roomZsetKey()))
}

func GetRoom(id string) (*Room, error) {
    roomRedisKey := roomKey(id)
    res, err := common.RedisClient.Multi(func(conn redis.Conn) {
        conn.Send("GET", roomRedisKey)
        conn.Send("TTL", roomRedisKey)
    })
    if err != nil {
        logs.Error("msg[获取房间失败] err[%s]", err.Error())
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    result := res.([]interface{})
    var room Room
    if result[0] == nil {
        return nil, nil
    }
    if err = common.DeJson(result[0].([]byte), &room); err != nil {
        logs.Error("msg[GetRoom解析room redis data err] err[%s]", err.Error())
        return nil, err
    }

    if ttl := result[1].(int64); ttl < 3600*3 {
        _, err = common.RedisClient.Expire(roomRedisKey, ROOM_LIVE_TIME)
        if err != nil {
            logs.Error("msg[延长room data redis 时间失败] err[%s]", err.Error())
        }
    }
    return &room, nil
}
func DelRoom(id string) (int, error) {
    res, err := common.RedisClient.Multi(func(conn redis.Conn) {
        conn.Send("DEL", roomUserKey(id))
        conn.Send("DEL", roomKey(id))
        conn.Send("ZREM", roomZsetKey(), id)
    })
    if err != nil {
        return 0, err
    }
    result := res.([]interface{})
    return redis.Int(result[1], nil)
}

func roomKey(id string) string {
    return "comet:room:" + id
}
func roomZsetKey() string {
    return "comet:room_zset"
}
func roomUserKey(roomId string) string {
    return "comet:roomUserList:" + roomId
}

func (r *Room) Users() (map[string]interface{}, error) {
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
        common.DeJson([]byte(st), &ru)
        res[dt] = ru
    }
    return res, nil
}

func (r *Room) Join(s *Session) (bool, error) {
    //r.users[s.Id] = RUser{SId:s.Id, Ip:s.IP, User:*s.User}
    //退出旧房间
    if s.RoomId != "" {
        beego.Debug("msg[用户%s退出旧房间%s]", s.DeviceToken, s.RoomId)
        oldR, _ := GetRoom(s.RoomId)
        if oldR != nil {
            oldR.Leave(s)
        }
    }
    s.RoomId = r.Id
    ru := RUser{DeviceToken: s.DeviceToken, User: *s.User, IP: s.IP}
    //查找当前用户是否已经在聊天室中
    user, err := common.RedisClient.Do("hget", roomUserKey(r.Id), ru.DeviceToken)
    if err != nil {
        return false, err
    }
    //用户已经在聊天室中直接返回成功
    if user != nil {
        return true, nil
    }
    jsonStr, err := common.EnJson(ru)
    if err != nil {
        return false, err
    }
    res, err := redis.Int(common.RedisClient.Do("hset", roomUserKey(r.Id), ru.DeviceToken, jsonStr))
    if err != nil {
        return false, err
    }
    if res < 1 {
        return false, errors.New("进入房间写入redis失败")
    }
    return true, nil
}

func (r *Room) Leave(s *Session) (bool, error) {
    if len(s.DeviceToken) < 1 {
        return false, errors.New("session.DeviceToken为空")
    }
    _, err := common.RedisClient.Do("hdel", roomUserKey(r.Id), s.DeviceToken)
    if err != nil {
        return false, err
    }
    s.RoomId = ""
    userNum, err := redis.Int(common.RedisClient.Do("hlen", roomUserKey(r.Id)))
    if err != nil {
        return false, err
    }

    if userNum < 1 {
        DelRoom(r.Id)
        beego.Debug("msg[房间%d内用户为空删除房间]", r.Id)
    }
    return true, nil
}

// This function handles all incoming chan messages.

//房间内广播
func (r *Room) Broadcast(msg *Msg) (bool, error) {
    resData := map[string]interface{}{}
    if msg.Data != "" {
        if err := common.DeJson([]byte(msg.Data), &resData); err != nil {
            logs.Error("msg[Broadcast DeJson err] err[%s]", err.Error())
            return false, err
        }
    } else {
        resData["room_id"] = r.Id
    }
    jsonByte, err := common.EnJson(resData)
    if err != nil {
        return false, err
    }
    msg.Data = string(jsonByte)
    msg.Type = TYPE_ROOM_MSG

    _, err = SaveRoomMsg(r.Id, msg)

    if err!=nil {
        logs.Error("msg[room msg 写入数据库失败]")
    }

    beego.Debug("msg[room broadcast]")
    users, err := r.Users()
    if err != nil {
        return false, err
    }
    //beego.Debug("msg[room中的用户] users[%s]", users)
    for _, user := range users {
        //session.Send(msg)
        tmsg := *msg
        tmsg.DeviceToken = user.(RUser).DeviceToken
        Server.Unicast(tmsg.DeviceToken, tmsg)
    }
    return true, nil
}

func SaveRoomMsg(roomId string, msg *Msg) (uint64, error){
    var msgData map[string]interface{}
    err := common.DeJson([]byte(msg.Data), &msgData)
    if err != nil {
        logs.Error("msg[SaveRoomMsg DeJson err] err[%s]", err.Error())
        return 0, err
    }
    content, ok := msgData["content"]
    if !ok {
        logs.Error("msg[SaveRoomMsg msg.Data 中不包含content]")
        return 0, errors.New("SaveRoomMsg msg.Data 中不包含content")
    }
    uid, ok := msgData["uid"]
    if !ok {
        logs.Error("msg[SaveRoomMsg msg.Data 中不包含uid]")
        return 0, errors.New("SaveRoomMsg msg.Data 中不包含content")
    }
    var uidInt int64
    if reflect.ValueOf(uid).Kind() == reflect.String{
        uidInt, _ = strconv.ParseInt(uid.(string), 10, 64)
    }else{
        uidInt = int64(uid.(float64))
    }
    roomMsg := &models.RoomMsg{RoomId: roomId, Content: content.(string), Uid: uidInt, CT: time.Now().Unix()}
    res := models.InsertRoomMsg(roomId, roomMsg)
    return roomMsg.Id, res.Error
}