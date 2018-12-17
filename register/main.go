package main

import (
	"container/list"
	"fmt"
	"net"
	"net/rpc"
)

var (
	ServiceMap = make(map[string]*list.Element)
	ServiceList = list.New()
)

type Manager struct {}

func (r *Manager) Register(ip string, res *int) error {
	if _, ok := ServiceMap[ip]; ok{
		*res = 0
		return nil
	}
	e := ServiceList.PushBack(ip)
	ServiceMap[ip] = e
	*res = 0
	fmt.Println(ServiceList)
	return nil
}

func (r *Manager) Unregister(ip string, res *int) error {
	if e, ok := ServiceMap[ip]; ok{
		ServiceList.Remove(e)
	}
	*res = 0
	return nil
}
func (r *Manager) List(req string, res *[]string) error{

	for item := ServiceList.Front(); item != nil; item = item.Next(){
		*res = append(*res, item.Value.(string))
		fmt.Println(item.Value.(string))
	}
	return nil
}

func main(){
	m := new(Manager)
	err := rpc.Register(m)
	if err!=nil {
		panic(err)
	}
	//ResolveTCPAddr返回TCP端点的地址。
	//网络必须是TCP网络名称。
	tcpAddr, err := net.ResolveTCPAddr("tcp",":1234")
	if err!=nil {
		panic(err)
	}
	listener,err := net.ListenTCP("tcp",tcpAddr)
	for  {
		//需要自己控制连接，当有客户端连接上来后，我们需要把这个连接交给rpc 来处理
		conn,err:=listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}
}