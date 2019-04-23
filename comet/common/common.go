package common

import (
    "crypto/md5"
    "encoding/hex"
    "net"
    "log"
    "github.com/dgryski/go-farm"
)

var localIp = ""
/**
 * 获取本地ip
 */
func GetLocalIp() string {

    if len(localIp) > 10 {
        return localIp
    }

    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Printf("msg[get_local_ip_failure] detail[%s]", err.Error())
        return ""
    }
    for _, address := range addrs {
        // check the address type and if it is not a loopback than display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                localIp = ipnet.IP.String()
            }
        }
    }
    return localIp
}

/**
 * 设成设备token
 */
func GenerateDeviceToken(deviceId string, appKey string) string {
    return Md5(deviceId + "|" + appKey)
}

/**
 * MD5编码
 */
func Md5(text string) string {
    encoder := md5.New()
    encoder.Write([]byte(text))
    return hex.EncodeToString(encoder.Sum(nil))
}

//对字符串的hash取模
func StrMod(str string, cardinal int) int{
    h := farm.Hash32([]byte(str))
    return int(h) % cardinal
}