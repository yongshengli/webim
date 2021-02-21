package server

import (
	"comet/common"
	"fmt"
	"net/rpc"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

type rpcClientPool chan *rpc.Client

//ServerContext Server context
type ServerContext struct {
	Server           *Server
	serverRpcPoolMap map[string]rpcClientPool
}

//Register 注册server
func (sm *ServerContext) Register(s Info) (int, error) {
	key := fmt.Sprintf("%s:%s", s.Host, s.Port)
	return redis.Int(common.RedisClient.Do("zadd", serverMapKey(), s.LastActive, key))
}

//List 列出全部server
func (sm *ServerContext) List() ([]Info, error) {
	//"WITHSCORES"
	serverArr, err := redis.Int64Map(common.RedisClient.Do("zrevrange", serverMapKey(), 0, -1, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	logs.Debug("msg[全部server列表] map[%v]", serverArr)
	timeNow := time.Now().Unix()
	res := make([]Info, 0)
	for s, at := range serverArr {
		tmp := strings.Split(s, ":")
		sInfo := Info{Host: tmp[0], Port: tmp[1], LastActive: at}
		leadTime := timeNow - at
		if leadTime > 60*3 { //客观下线
			if leadTime > 600*5 { //物理下线
				sm.Remove(sInfo)
			}
			continue
		}
		res = append(res, sInfo)
	}
	logs.Debug("msg[可用server列表] map[%v]", res)
	return res, nil
}

//Len 统计server个数
func (sm *ServerContext) Len() (int, error) {
	return redis.Int(common.RedisClient.Do("zcard", serverMapKey()))
}

//Remove 移除server
func (sm *ServerContext) Remove(s Info) (int, error) {
	key := fmt.Sprintf("%s:%s", s.Host, s.Port)
	logs.Info("remove server " + key)
	return redis.Int(common.RedisClient.Do("zrem", serverMapKey(), key))
}

func serverMapKey() string {
	return "comet:serverMap"
}
