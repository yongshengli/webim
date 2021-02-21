package server

import (
	"comet/common"
	"comet/server/base"
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
)

func delDeviceTokenInfo(deviceToken string) (int, error) {
	logs.Debug("msg[call_delDeviceTokenInfo] device_toke[%s]", deviceToken)
	return common.RedisClient.Del([]string{deviceTokenKey(deviceToken)})
}
func saveDeviceTokenInfo(user *base.User) (string, error) {
	if len(user.DeviceToken) < 1 {
		return "", errors.New("DeviceToken为空")
	}

	if jsonStr, err := common.EnJson(user); err != nil {
		beego.Error(err)
		return "", err
	} else {
		return common.RedisClient.Set(deviceTokenKey(user.DeviceToken), jsonStr, base.SESSION_LIVE_TIME)
	}
}
func getDeviceTokenByUid(uid string) (string, error) {
	return redis.String(common.RedisClient.Get(uidKey(uid)))
}

func getDeviceTokenInfo(deviceToken string) (*base.User, error) {
	user := &base.User{}
	tokenKey := deviceTokenKey(deviceToken)
	commands := make([]common.RedisCommands, 2)
	commands[0] = common.RedisCommands{CommandName: "GET", Args: []interface{}{tokenKey}}
	commands[1] = common.RedisCommands{CommandName: "TTL", Args: []interface{}{tokenKey}}
	reply := common.RedisClient.Pipeline(commands)
	for i := 0; i < len(commands); i++ {
		if reply[i]["err"] != nil {
			return nil, reply[i]["err"].(error)
		}
	}
	if reply[0]["reply"] == nil {
		return nil, nil
	}
	if err := common.DeJson(reply[0]["reply"].([]byte), user); err != nil {
		logs.Error("msg[解析json失败] method[getDeviceTokenInfo] err[%s]", err.Error())
		return nil, err
	}
	if ttl := reply[1]["reply"].(int64); ttl < 3600 {
		if _, err := common.RedisClient.Expire(tokenKey, base.SESSION_LIVE_TIME); err != nil {
			logs.Error("msg[延长session有效期失败] err[%s]", err.Error())
		}
	}
	return user, nil
}
func deviceTokenKey(deviceToken string) string {
	return "comet:token:" + deviceToken
}

func uidKey(uid string) string {
	return "comet:uid:" + uid
}
