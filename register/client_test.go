package main

import (
	"fmt"
	"net/rpc"
	"testing"
)

func TestClient(t *testing.T){
	client,err:=rpc.Dial("tcp","127.0.0.1:1234")
	if err!=nil {
		t.Fatal("连接Dial的发生了错误，我要退出了",err)
	}
	var res int
	err2 := client.Call("Manager.Register", "127.0.0.1:8000", &res)

	if err2!=nil {
		t.Fatal(err)
	}
	fmt.Println(res)
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