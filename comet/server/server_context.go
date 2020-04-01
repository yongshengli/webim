package server

import (
	"fmt"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"time"
	"webim/comet/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

const RPC_CHAN_NUM = 30

type rpcClinePool chan *rpc.Client
type Context struct {
	server           *server
	serverRpcPoolMap map[string]rpcClinePool
}

func (sm *Context) CallRpcClient(host string, method, args string, reply interface{}) error {
	if rpcPool, ok := sm.serverRpcPoolMap[host]; ok {
		client := <-rpcPool
		return client.Call(method, args, reply)
	}
	return nil
}

func (sm *Context) createRpcClinePool() {
	serverMap, err := sm.List()
	if err != nil {
		logs.Error("msg[createRpcClientPool err] err[%s]", err.Error())
	}
	for _, s := range serverMap {
		sm.addServerRpcClinePool(s.Host, s.Port)
	}
}

func (sm *Context) addServerRpcClinePool(host, port string) bool {
	if host == sm.server.Host {
		return false
	}
	if _, ok := sm.serverRpcPoolMap[host]; ok {
		return true
	}
	sm.serverRpcPoolMap[host] = make(chan *rpc.Client, RPC_CHAN_NUM)
	addr := host + ":" + port
	for i := 0; i < 10; i++ {
		client, err := jsonrpc.Dial("tcp", addr)
		if err != nil {
			logs.Error("msg[newServerRpcChan err] err[%s]", err.Error())
			continue
		}
		sm.serverRpcPoolMap[host] <- client
	}
	return true
}

//Register 注册server
func (sm *Context) Register(s Info) (int, error) {
	key := fmt.Sprintf("%s:%d", s.Host, s.Port)
	return redis.Int(common.RedisClient.Do("zadd", serverMapKey(), s.LastActive, key))
}

//List 列出全部server
func (sm *Context) List() ([]Info, error) {
	//"WITHSCORES"
	serverArr, err := redis.Int64Map(common.RedisClient.Do("zrevrange", serverMapKey(), 0, -1))
	if err != nil {
		return nil, err
	}

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
	return res, nil
}

//Len 统计server个数
func (sm *Context) Len() (int, error) {
	return redis.Int(common.RedisClient.Do("HLEN", serverMapKey()))
}

//Remove 移除server
func (sm *Context) Remove(s Info) (int, error) {
	key := fmt.Sprintf("%s:%d", s.Host, s.Port)
	beego.Info("remove server " + key)
	return redis.Int(common.RedisClient.Do("zrem", serverMapKey(), key))
}

func serverMapKey() string {
	return "comet:serverMap"
}
