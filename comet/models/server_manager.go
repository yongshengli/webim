package models

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
)

type serverManger func()

var ServerManager = new(serverManger)

func (sm *serverManger) Register(addr string) (int, error){
    return redis.Int(common.RedisClient.Do("hset", serverMapKey(), addr, addr))
}

func (sm *serverManger) List() (map[string]string, error){
    replay, err := common.RedisClient.Do("hgetall", serverMapKey())
    if err!=nil{
        return nil, err
    }
    if replay==nil{
        return nil, nil
    }
    return redis.StringMap(replay, err)
}

func (sm *serverManger) Remove(addr string) (int, error){
    return redis.Int(common.RedisClient.Do("hdel", serverMapKey(), []string{addr}))
}

func serverMapKey() string{
    return "comet:serverMap"
}