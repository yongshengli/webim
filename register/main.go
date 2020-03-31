package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"webim/comet/server"
)

type Manager struct {
	ServiceMap map[*net.Conn]string
}

var Server = &Manager{
	ServiceMap: make(map[*net.Conn]string),
}

func (m *Manager) Register(conn *net.Conn) error {
	if _, ok := m.ServiceMap[conn]; ok {
		return nil
	}
	m.ServiceMap[conn] = (*conn).RemoteAddr().String()
	log.Println(m.ServiceMap)
	return nil
}

func (m *Manager) Unregister(conn *net.Conn) error {
	log.Println("用户")
	if _, ok := m.ServiceMap[conn]; ok {
		delete(m.ServiceMap, conn)
	}
	return nil
}
func (m *Manager) List() []string {
	var slice []string
	for _, v := range m.ServiceMap {
		slice = append(slice, v)
		fmt.Println(v)
	}
	return slice
}
func (m *Manager) broadcast(msg server.Msg) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for conn := range m.ServiceMap {
		_, err2 := (*conn).Write(b)
		if err2 != nil {
			continue
		}
	}
	return nil
}
func handleConnection(conn *net.Conn) {
	defer func(conn *net.Conn) {
		defer (*conn).Close()
		Server.Unregister(conn)
	}(conn)
	Server.Register(conn)

	for {
		var buf [512]byte
		n, err := (*conn).Read(buf[0:])
		if n > 0 {
			if string(buf[:n]) == "list" {
				b, _ := json.Marshal(Server.List())
				_, err2 := (*conn).Write(b)
				if err2 != nil {
					return
				}
			}
		}
		if err != nil {
			return
		}
	}
}
func main() {
	//ResolveTCPAddr返回TCP端点的地址。
	//网络必须是TCP网络名称。
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		//需要自己控制连接，当有客户端连接上来后，我们需要把这个连接交给rpc 来处理
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(&conn)
	}
}
