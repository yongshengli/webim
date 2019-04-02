package models

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
    "encoding/json"
    "github.com/astaxie/beego"
)

type serverManger func()

type ServerInfo struct {
    Host string `json:"host"`
    Port string `json:"port"`
    Data map[string]string `json:"data"`
}
var (
    ServerManager = new(serverManger)
    CurrentServer ServerInfo
)

func (sm *serverManger) Register(port string) (int, error){
    CurrentServer = ServerInfo{Host:common.GetLocalIp(), Port:port}
    b, err := json.Marshal(CurrentServer)
    if err!=nil{
        beego.Error(err)
        return 0, err
    }
    return redis.Int(common.RedisClient.Do("hset", serverMapKey(), CurrentServer.Host, string(b)))
}

func (sm *serverManger) List() (map[string]ServerInfo, error) {
    replay, err := common.RedisClient.Do("hgetall", serverMapKey())
    if err != nil {
        return nil, err
    }
    if replay == nil {
        return nil, nil
    }
    strM, err := redis.StringMap(replay, err)
    if err!=nil{
        beego.Error(err)
    }
    res := make(map[string]ServerInfo, len(strM))
    for h, v := range strM{
        t := ServerInfo{}
        json.Unmarshal([]byte(v), &t)
        res[h] = t
    }
    return res, nil
}

func (sm *serverManger) Remove() (int, error){
    beego.Info("remove server " + CurrentServer.Host)
    return redis.Int(common.RedisClient.Do("hdel", serverMapKey(), []string{CurrentServer.Host}))
}

func serverMapKey() string{
    return "comet:serverMap"
}