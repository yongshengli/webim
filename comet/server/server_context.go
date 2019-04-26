package server

import (
    "webim/comet/common"
    "github.com/gomodule/redigo/redis"
    "encoding/json"
    "github.com/astaxie/beego"
    "time"
    "net/rpc"
    "net/rpc/jsonrpc"
    "github.com/astaxie/beego/logs"
)
const RPC_CHAN_NUM = 30

type rpcPool chan *rpc.Client
type Context struct {
    server server
    serverRpcPoolMap map[string]rpcPool
}

func (sm *Context) CallRpcClient(host string, method, args string, reply interface{}) error{
    if rpcPool, ok := sm.serverRpcPoolMap[host]; ok {
        client := <-rpcPool
        return client.Call(method, args, reply)
    }
    return nil
}
func (sm *Context) addServerRpcPool(host, port string) bool{
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
        if err!=nil {
            logs.Error("msg[newServerRpcChan err] err[%s]", err.Error())
            continue
        }
        sm.serverRpcPoolMap[host] <- client
    }
    return true
}
func (sm *Context) Register(host string, server Info) (int, error){
    server.LastActive = time.Now().Unix()
    b, err := json.Marshal(server)
    if err != nil {
        beego.Error(err)
        return 0, err
    }
    return redis.Int(common.RedisClient.Do("HSET", serverMapKey(), host, string(b)))
}

func (sm *Context) List() (map[string]Info, error) {
    replay, err := common.RedisClient.Do("HGETALL", serverMapKey())
    if err != nil {
        return nil, err
    }
    if replay == nil {
        return nil, nil
    }
    strM, err := redis.StringMap(replay, err)
    if err != nil {
        beego.Error(err)
    }
    timeNow :=  time.Now().Unix()
    res := make(map[string]Info, len(strM))
    for h, v := range strM {
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

func (sm *Context) Len() (int, error) {
    return redis.Int(common.RedisClient.Do("HLEN", serverMapKey()))
}

func (sm *Context) Remove(host string) (int, error){
    beego.Info("remove server " + host)
    return redis.Int(common.RedisClient.Do("HDEL", serverMapKey(), host))
}

func serverMapKey() string{
    return "comet:serverMap"
}