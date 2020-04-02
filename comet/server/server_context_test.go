package server

import (
	"fmt"
	"testing"
	"time"
	"webim/comet/common"

	"github.com/gomodule/redigo/redis"
)

func TestServerRegister(t *testing.T) {
	server := Info{
		Host:       "127.0.0.1",
		Port:       "8000",
		LastActive: time.Now().Unix(),
	}
	context := new(Context)

	if _, err := context.Register(server); err != nil {
		t.Error(err)
	}
	//fmt.Println(res)
	list, err := context.List()
	if err != nil {
		t.Error(err)
	}
	if len(list) < 1 {
		t.Error("没有取到主机map")
	}
	fmt.Println(list)
	context.Remove(server)
	redisKey := fmt.Sprintf("%s:%s", server.Host, server.Port)
	num, err := redis.Int64(common.RedisClient.Do("ZSCORE", serverMapKey(), redisKey))
	if err != nil && err != redis.ErrNil {
		t.Error(err)
	}
	if num > 0 {
		t.Error("删除失败")
	}
}
func init() {
	common.RedisInitTest()
}
