package models

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
)

type serverManger func()

var (
    ServerManager = new(serverManger)
)

func (sm *serverManger) Register(port string) (int, error){
    return redis.Int(common.RedisClient.Do("hset", serverMapKey(), common.GetLocalIp(), port))
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

func (sm *serverManger) Remove(ip string) (int, error){
    return redis.Int(common.RedisClient.Do("hdel", serverMapKey(), []string{ip}))
}

func serverMapKey() string{
    return "comet:serverMap"
}