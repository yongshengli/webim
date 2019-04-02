package models

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
)

type serverManger func()

type ServerInfo struct {
    Host string
    Port string
    Data map[string]string
}
var (
    ServerManager = new(serverManger)
    CurrentServer ServerInfo
)

func (sm *serverManger) Register(port string) (int, error){
    addr := common.GetLocalIp()+":"+port
    CurrentServer = ServerInfo{Host:common.GetLocalIp(), Port:port}

    return redis.Int(common.RedisClient.Do("hset", serverMapKey(), addr, CurrentServer))
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