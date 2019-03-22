package common

import (
    "errors"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "github.com/gomodule/redigo/redis"
    "sync"
    "time"
)

var (
    once   = sync.Once{}
    RedisClient *redisClient
)

type redisClient struct {
    pool *redis.Pool
}

func (r *redisClient) Get(key string) (interface{}, error) {
    reply, err := r.Do("GET", key)
    if err != nil {
        logs.Error("msg[redis_get_failed:%s]", err.Error())
    }
    return reply, err
}

func (r *redisClient) Ttl(key string) (int64, error) {
    ttl, err := r.Do("ttl", key)
    if err != nil {
        logs.Error("msg[redis_get_failed:%s]", err.Error())
    }
    return ttl.(int64), err
}

func (r *redisClient) GetString(key string) string {
    rep, _ := redis.String(r.Get(key))
    return rep
}
func (r *redisClient) Set(key string, val interface{}, timeout time.Duration) string {
    reply, err := redis.String(r.Do("SETEX", key, int64(timeout/time.Second), val))
    if err != nil {
        logs.Error("msg[redis_setex_failed:%s]", err.Error())
        return ""
    }
    return reply
}
func (r *redisClient) MgetString(keys []string) []string {
    tKeys := arrStr2ArrInterface(keys)
    rep, err := redis.Strings(r.Do("MGET", tKeys...))
    if err != nil {
        logs.Error("msg[redis_mget_failed:%s]", err.Error())
    }
    return rep
}
func (r *redisClient) Mget(keys []string) ([]interface{}, error) {
    tKeys := arrStr2ArrInterface(keys)
    rep, err := redis.Values(r.Do("MGET", tKeys...))
    if err != nil {
        logs.Error("msg[redis_mget_failed:%s]", err.Error())
    }
    return rep, err
}

func (r *redisClient) Exists(key string) bool{
    rep, err :=redis.Bool(r.Do("EXISTS", key))
    if err!=nil{
        logs.Error("msg[redis_exists_failed:%s]", err.Error())
        return false
    }
    return rep
}

func (r *redisClient) Del(keys []string) int {
    tKeys := arrStr2ArrInterface(keys)
    num, err := redis.Int(r.Do("DEL", tKeys...))
    if err != nil {
        logs.Error("msg[redis_del_failed:%s]", err.Error())
        return 0
    }
    return num
}
// Incr increase counter in redis.
func (r *redisClient) Incr(key string) bool {
    rep, err := redis.Bool(r.Do("INCR", key))
    if err != nil {
        logs.Error("msg[redis_incr_failed:%s]", err.Error())
        return false
    }
    return rep
}

// Decr decrease counter in redis.
func (r *redisClient) Decr(key string) bool {
    rep, err := redis.Bool(r.Do("DECR", key))
    if err != nil {
        logs.Error("msg[redis_decr_failed:%s]", err.Error())
        return false
    }
    return rep
}
func (r *redisClient) initRedisPool() {
    once.Do(func() {
        dialFunc := func() (c redis.Conn, err error) {
            addr := beego.AppConfig.String("redis.host") + ":" + beego.AppConfig.String("redis.port")
            c, err = redis.Dial("tcp", addr)
            if err != nil {
                logs.Info("msg[conn_to_redis_failure:%s] addr[%s]", err.Error(), addr)
                return nil, err
            }
            _, selecterr := c.Do("SELECT", 0)
            if selecterr != nil {
                c.Close()
                return nil, selecterr
            }
            logs.Info("msg[get_redis_conn_success] addr[%s]", addr)
            return
        }
        // initialize a new pool
        r.pool = &redis.Pool{
            Wait:true,//，当程序执行get()，无法获得可用连接时，将会暂时阻塞。
            MaxIdle:     500,
            MaxActive:	 5000,//设MaxActive=0(表示无限大)或者足够大。
            IdleTimeout: 180 * time.Second,
            Dial:        dialFunc,
        }
    })
}

func arrStr2ArrInterface(keys []string) []interface{} {
    tKeys := make([]interface{}, len(keys))
    for _, key := range keys {
        tKeys = append(tKeys, key)
    }
    return tKeys
}
func (r *redisClient) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
    if len(args) < 1 {
        return nil, errors.New("missing required arguments")
    }
    conn := r.pool.Get()
    defer conn.Close()
    return conn.Do(commandName, args...)
}

func init() {
    RedisClient = &redisClient{}
    RedisClient.initRedisPool()
}