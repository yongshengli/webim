package server

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

//RpcService Rpc Service
type RpcService struct {
	s *server
}

//Unicast 单播
func (rf *RpcService) Unicast(args map[string]interface{}, reply *bool) error {
	if _, ok := args["device_token"]; !ok {
		return errors.New("device_token不能为空")
	}
	if _, ok := args["msg"]; !ok {
		return errors.New("msg不能为空")
	}
	msg, err := Map2Msg(args["msg"].(map[string]interface{}))
	if err != nil {
		return err
	}
	res, err := rf.s.Unicast(args["device_token"].(string), msg)
	*reply = res
	return err
}

//Broadcast 广播
func (rf *RpcService) Broadcast(args map[string]interface{}, reply *bool) error {
	if _, ok := args["msg"]; !ok {
		return errors.New("msg不能为空")
	}
	if _, ok := args["msg"].(map[string]interface{}); !ok {
		return errors.New("msg格式错误")
	}
	msg, err := Map2Msg(args["msg"].(map[string]interface{}))
	if err != nil {
		return err
	}
	res, err := rf.s.Broadcast(msg)
	*reply = res
	return err
}

//BroadcastSelf 只广播本机
func (rf *RpcService) BroadcastSelf(args map[string]interface{}, reply *bool) error {
	if _, ok := args["msg"]; !ok {
		return errors.New("msg不能为空")
	}
	if _, ok := args["msg"].(map[string]interface{}); !ok {
		return errors.New("msg格式错误")
	}
	msg, err := Map2Msg(args["msg"].(map[string]interface{}))
	if err != nil {
		return err
	}
	res, err := rf.s.BroadcastSelf(msg)
	*reply = res
	return err
}

//Status 状态
func (rf *RpcService) Status(args map[string]interface{}, reply *map[string]interface{}) error {

	return nil
}

//Ping ping
func (rf *RpcService) Ping(args map[string]interface{}, reply *string) error {
	*reply = "pong"
	return nil
}

//RunRpcService 执行
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
