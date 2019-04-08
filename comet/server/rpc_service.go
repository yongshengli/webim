package server

import (
    "net/rpc"
    "net"
    "fmt"
    "os"
    "net/rpc/jsonrpc"
    "errors"
)

type RpcService struct {
    s *server
}

func (rf *RpcService) Unicast(args map[string]interface{}, reply *bool) error {
    if _, ok := args["device_token"]; !ok {
        return errors.New("device_token不能为空")
    }
    if _, ok := args["msg"]; !ok {
        return errors.New("msg不能为空")
    }
    if _, tok := args["msg"].(Msg); !tok {
        return errors.New("msg格式错误")
    }
    res, err := rf.s.Unicast(args["device_token"].(string), args["msg"].(Msg))
    *reply = res
    return err
}

func (rf *RpcService) Broadcast(args map[string]interface{}, reply *bool) error {
    if _, ok := args["msg"]; !ok {
        return errors.New("msg不能为空")
    }
    if _, tok := args["msg"].(Msg); !tok {
        return errors.New("msg格式错误")
    }
    res, err := rf.s.Broadcast(args["msg"].(Msg))
    *reply = res
    return err
}

//只广播本机
func (rf *RpcService) BroadcastSelf(args map[string]interface{}, reply *bool) error {
    if _, ok := args["msg"]; !ok {
        return errors.New("msg不能为空")
    }
    if _, tok := args["msg"].(Msg); !tok {
        return errors.New("msg格式错误")
    }
    res, err := rf.s.BroadcastSelf(args["msg"].(Msg))
    *reply = res
    return err
}

func (rf *RpcService) Status(args map[string]interface{}, reply *map[string]interface{}) error {

    return nil
}

func (rf *RpcService) Ping(args map[string]interface{}, reply *string) error {
    *reply = "pong"
    return nil
}

func RunRpcService(s *server) {
    rpc.Register(&RpcService{s: s})
    tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+s.Port)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    listener, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for {
        // todo   需要自己控制连接，当有客户端连接上来后，我们需要把这个连接交给rpc 来处理
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        jsonrpc.ServeConn(conn)
    }
}