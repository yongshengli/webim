package common

import (
    "net"
    "log"
)

func GetLocalIp() (ipAddr string){
    ipAddr = ""
    addrSlice, err := net.InterfaceAddrs()
    if nil != err {
        log.Println("Get local IP addr failed!!!")
        return
    }
    for _, addr := range addrSlice {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if nil != ipnet.IP.To4() {
                ipAddr = ipnet.IP.String()
                return
            }
        }
    }
    return
}