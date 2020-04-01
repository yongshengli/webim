package common

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"net"

	"github.com/dgryski/go-farm"
	guuid "github.com/satori/go.uuid"
)

//GetLocalIp 获取本地ip
func GetLocalIp() string {
	var localIp = ""
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

//GenerateDeviceToken 设成设备token
func GenerateDeviceToken(deviceId string, appKey string) string {
	return Md5(deviceId + "|" + appKey)
}

//Md5 字符串MD5编码
func Md5(text string) string {
	encoder := md5.New()
	encoder.Write([]byte(text))
	return hex.EncodeToString(encoder.Sum(nil))
}

//StrMod 对字符串的hash取模
func StrMod(str string, cardinal int) int {
	h := farm.Hash32([]byte(str))
	return int(h) % cardinal
}

//Map2String 将map 装换为json string
func Map2String(m map[string]interface{}) (string, error) {
	resByte, err := EnJson(m)
	if err != nil {
		return "", err
	}
	return string(resByte), nil
}

//Uuid 获取UUID
func Uuid() string {
	return guuid.NewV4().String()
}

//LogId 获取logid
func LogId() string {
	return guuid.NewV1().String()
}
