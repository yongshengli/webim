package models

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
)

type serverManger func()

var (
    ServerManager = new(serverManger)
    CurrentServer = map[string]string{}
)

func (sm *serverManger) Register(port string) (int, error){
    addr := common.GetLocalIp()+":"+port
    CurrentServer["host"] = common.GetLocalIp()
    CurrentServer["port"] = port

    return redis.Int(common.RedisClient.Do("hset", serverMapKey(), addr, port))
}

func (sm *serverManger) List() (map[string]string, error) {
    replay, err := common.RedisClient.Do("hgetall", serverMapKey())
    if err != nil {
        return nil, err
    }
    if replay == nil {
        return nil, nil
    }
    return redis.StringMap(replay, err)
}

func (sm *serverManger) Remove() (int, error){
    return redis.Int(common.RedisClient.Do("hdel", serverMapKey(), []string{common.GetLocalIp()}))
}

func serverMapKey() string{
    return "comet:serverMap"
}