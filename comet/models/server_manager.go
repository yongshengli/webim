package models

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
    "encoding/json"
    "github.com/astaxie/beego"
    "time"
    "github.com/astaxie/beego/logs"
)

type serverManger func()

type ServerInfo struct {
    Host       string            `json:"host"`
    Port       string            `json:"port"`
    Data       map[string]string `json:"data"`
    LastActive int64             `json:"last_active"` //上次活跃时间
}

var (
    ServerManager = new(serverManger)
    CurrentServer ServerInfo
)

func (sm *serverManger) ReportLive(port string){
    sm.Register(port)

    t := time.NewTicker(time.Minute)
    defer t.Stop()
    for {
        <- t.C
        sm.Register(port)
        logs.Debug("msg[服务报活]")
    }
}

func (sm *serverManger) Register(port string) (int, error){
    CurrentServer = ServerInfo{Host:common.GetLocalIp(), Port:port, LastActive:time.Now().Unix()}
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
    timeNow :=  time.Now().Unix()
    res := make(map[string]ServerInfo, len(strM))
    for h, v := range strM{
        t := ServerInfo{}
        json.Unmarshal([]byte(v), &t)
        leadTime := timeNow - t.LastActive
        if leadTime > 600*3{ //客观下线
            if leadTime > 600*5{ //物理下线
                sm.Remove(t.Host)
            }
            continue
        }
        res[h] = t
    }
    return res, nil
}

func (sm *serverManger) Remove(host string) (int, error){
    beego.Info("remove server " + host)
    return redis.Int(common.RedisClient.Do("hdel", serverMapKey(), host))
}

func serverMapKey() string{
    return "comet:serverMap"
}