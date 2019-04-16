package server

import (
    "github.com/astaxie/beego/logs"
    "webim/comet/common"
    "github.com/astaxie/beego"
    "github.com/gomodule/redigo/redis"
    "errors"
)

func delDeviceTokenInfo(deviceToken string) (int, error) {
    logs.Debug("msg[call_delDeviceTokenInfo] device_toke[%s]", deviceToken)
    return common.RedisClient.Del([]string{deviceTokenKey(deviceToken)})
}
func saveDeviceTokenInfo(user *User) (string, error) {
    if len(user.DeviceToken) < 1 {
        return "", errors.New("DeviceToken为空")
    }

    if jsonStr, err := common.EnJson(user); err != nil {
        beego.Error(err)
        return "", err
    }else {
        return common.RedisClient.Set(deviceTokenKey(user.DeviceToken), jsonStr, SESSION_LIVE_TIME)
    }
}
func getDeviceTokenByUid(uid string) (string, error) {
    return redis.String(common.RedisClient.Get(uidKey(uid)))
}

func getDeviceTokenInfo(deviceToken string) (map[string]string, error) {
    var res map[string]string
    tokenKey := deviceTokenKey(deviceToken)
    replay, err := common.RedisClient.Multi(func (conn redis.Conn){
        conn.Send("GET", tokenKey)
        conn.Send("TTL", tokenKey)
    })
    if replay == nil {
        return nil, err
    }
    tmpRes := replay.([]interface{})
    if tmpRes[0] == nil {
        return nil, nil
    }
    if err = common.DeJson(tmpRes[0].([]byte), &res); err != nil {
        logs.Error("msg[解析json失败] method[getDeviceTokenInfo] err[%s]", err.Error())
        return nil, err
    }
    if ttl := tmpRes[1].(int64); ttl < 3600 {
        _, err = common.RedisClient.Expire(tokenKey, SESSION_LIVE_TIME)
        if err != nil {
            logs.Error("msg[延长session有效期失败] err[%s]", err.Error())
        }
    }
    return res, nil
}
func deviceTokenKey(deviceToken string) string {
    return "comet:token:" + deviceToken
}

func uidKey(uid string) string {
    return "comet:uid:" + uid
}