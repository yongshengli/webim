package server

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
    "encoding/json"
    "github.com/astaxie/beego"
    "time"
)

type Context struct {
}

func (sm *Context) Register(host string, server Info) (int, error){
    server.LastActive = time.Now().Unix()
    b, err := json.Marshal(server)
    if err != nil {
        beego.Error(err)
        return 0, err
    }
    return redis.Int(common.RedisClient.Do("hset", serverMapKey(), host, string(b)))
}

func (sm *Context) List() (map[string]Info, error) {
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
    res := make(map[string]Info, len(strM))
    for h, v := range strM{
        t := Info{}
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

func (sm *Context) Remove(host string) (int, error){
    beego.Info("remove server " + host)
    return redis.Int(common.RedisClient.Do("hdel", serverMapKey(), host))
}

func serverMapKey() string{
    return "comet:serverMap"
}