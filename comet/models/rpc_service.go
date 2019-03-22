package models

import (
    "net/rpc"
    "net"
    "fmt"
    "os"
    "webim/comet/models/room"
)

type RpcFunc func()

func (rs *RpcFunc) Unicast(args map[string]interface{}, reply *bool) error {
    res, err := room.SessionManager.SendMsg(args["sid"].(string), args["msg"].(Msg))
    *reply = res
    return err
}

func (rs *RpcFunc) Broadcast(args map[string]interface{}, reply *bool) error {
    res, err := room.SessionManager.BroadcastSelf(args["msg"].(Msg))
    *reply = res
    return err
}

func RunRpcServer(port string) {
    rpcFunc := new(RpcFunc)
    rpc.Register(rpcFunc)
    tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+port)

    if err != nil {
        fmt.Println("错误了哦")
        os.Exit(1)
    }
    listener, err := net.ListenTCP("tcp", tcpAddr)
    for {
        // todo   需要自己控制连接，当有客户端连接上来后，我们需要把这个连接交给rpc 来处理
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        rpc.ServeConn(conn)
    }
}
