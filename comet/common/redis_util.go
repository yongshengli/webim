package common

import (
	"errors"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

var (
	RedisClient = &redisClient{}
)

type redisClient struct {
	sync.RWMutex // only used outsite for bootStrap
	pool         *redis.Pool
	done         bool
}

func (r *redisClient) Get(key string) (interface{}, error) {
	reply, err := r.Do("GET", key)
	if err != nil {
		logs.Error("msg[redis get err] err[%s] key[%s]", err.Error(), key)
	}
	return reply, err
}

func (r *redisClient) Ttl(key string) (int64, error) {
	ttl, err := r.Do("ttl", key)
	if err != nil {
		logs.Error("msg[redis ttl err] err[%s] key[%s]", err.Error(), key)
	}
	return redis.Int64(ttl, err)
}

func (r *redisClient) Expire(key string, timeout time.Duration) (int64, error) {
	res, err := r.Do("EXPIRE", key, int64(timeout/time.Second))
	if err != nil {
		logs.Error("msg[redis Expire err] err[%s] key[%s]", err.Error(), key)
	}
	return redis.Int64(res, err)
}

func (r *redisClient) GetString(key string) string {
	rep, _ := redis.String(r.Get(key))
	return rep
}
func (r *redisClient) Set(key string, val interface{}, timeout time.Duration) (string, error) {
	reply, err := redis.String(r.Do("SETEX", key, int64(timeout/time.Second), val))
	if err != nil {
		logs.Error("msg[redis setex err] err[%s] key[%s] val[%v]", err.Error(), key, val)
		return "", err
	}
	return reply, nil
}
func (r *redisClient) MgetString(keys []string) ([]string, error) {
	tKeys := arrStr2ArrInterface(keys)
	rep, err := redis.Strings(r.Do("MGET", tKeys...))
	if err != nil {
		logs.Error("msg[redis mget string err] err[%s]", err.Error())
	}
	return rep, err
}
func (r *redisClient) Mget(keys []string) ([]interface{}, error) {
	tKeys := arrStr2ArrInterface(keys)
	rep, err := redis.Values(r.Do("MGET", tKeys...))
	if err != nil {
		logs.Error("msg[redis mget err] err[%s]", err.Error())
	}
	return rep, err
}

func (r *redisClient) Exists(key string) (bool, error) {
	rep, err := redis.Bool(r.Do("EXISTS", key))
	if err != nil {
		logs.Error("msg[redis exists err] err[%s] key[%s]", err.Error(), key)
	}
	return rep, err
}

func (r *redisClient) Del(keys []string) (int, error) {
	tKeys := arrStr2ArrInterface(keys)
	num, err := redis.Int(r.Do("DEL", tKeys...))
	if err != nil {
		logs.Error("msg[redis del err] err[%s]", err.Error())
	}
	return num, err
}

// Incr increase counter in redis.
func (r *redisClient) Incr(key string) (uint64, error) {
	rep, err := redis.Uint64(r.Do("INCR", key))
	if err != nil {
		logs.Error("msg[redis incr err] err[%s] key[%s]", err.Error(), key)
	}
	return rep, err
}

// Decr decrease counter in redis.
func (r *redisClient) Decr(key string) (uint64, error) {
	rep, err := redis.Uint64(r.Do("DECR", key))
	if err != nil {
		logs.Error("msg[redis decr err] err[%s] key[%s]", err.Error(), key)
	}
	return rep, err
}
func (r *redisClient) initRedisPool(conf map[string]string) {
	dialFunc := func() (c redis.Conn, err error) {
		addr := conf["host"] + ":" + conf["port"]
		c, err = redis.Dial("tcp", addr)
		if err != nil {
			logs.Info("msg[conn_to_redis_failure] err[%s] addr[%s]", err.Error(), addr)
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
		Wait:        true, //，当程序执行get()，无法获得可用连接时，将会暂时阻塞。
		MaxIdle:     500,
		MaxActive:   5000, //设MaxActive=0(表示无限大)或者足够大。
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
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

func (r *redisClient) GetConn() redis.Conn {
	return r.pool.Get()
}

type RedisCommands struct {
	CommandName string
	Args        []interface{}
}

func (r *redisClient) Pipeline(commands []RedisCommands) []map[string]interface{} {
	conn := r.pool.Get()
	defer conn.Close()
	commandsNum := len(commands)
	for i := 0; i < commandsNum; i++ {
		conn.Send(commands[i].CommandName, commands[i].Args...)
	}
	conn.Flush()
	res := make([]map[string]interface{}, commandsNum)
	for i := 0; i < commandsNum; i++ {
		reply, terr := conn.Receive()
		res[i] = map[string]interface{}{"reply": reply, "err": terr}
	}
	return res
}

func (r *redisClient) Multi(callback func(conn redis.Conn)) (reply interface{}, err error) {
	conn := r.pool.Get()
	defer conn.Close()
	conn.Send("MULTI")
	callback(conn)
	return conn.Do("EXEC")
}

func RedisInit(conf map[string]string) {
	if RedisClient.done {
		return
	}
	RedisClient.Lock()
	defer RedisClient.Unlock()
	RedisClient.initRedisPool(conf)
}
func RedisInitTest() {
	RedisInit(map[string]string{
		"host": "127.0.0.1",
		"port": "6379",
	})
}
