package main

import (
	"net/rpc"
	"testing"
	"encoding/json"
	"net"
	"fmt"
)

func TestClient(t *testing.T){
	client, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		t.Fatal("连接Dial的发生了错误，我要退出了", err)
	}
	defer client.Close()
	buffer := make([]byte, 512)
	n, err := client.Write([]byte("list"))
	if err != nil {
		panic(err)
	}
	n, err2 := client.Read(buffer)

	var slice []string
	json.Unmarshal(buffer, slice)
	if err2 != nil {
		fmt.Println("Read failed:", err2)
		return
	}
	fmt.Println("count:", n, "msg:", string(buffer))
}
func TestList(t *testing.T){
	client,err:=rpc.Dial("tcp","127.0.0.1:1234")
	if err!=nil {
		t.Fatal(err)
	}
	var res []string
	err2 := client.Call("Manager.List", "", &res)

	if err2!=nil {
		t.Fatal(err)
	}
	for _, val := range res{
		if val!="127.0.0.1:8000"{
			t.Errorf("期望127.0.0.1:8000，but get %s", val)
		}
	}
}