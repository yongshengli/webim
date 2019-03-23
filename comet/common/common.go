package common

import (
    "net"
    "log"
)

var localIp = ""

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