package main

import (
    "net"
    "encoding/json"
    "fmt"
)

func main(){
    client, err := net.Dial("tcp", "127.0.0.1:1234")
    if err != nil {
        panic(err)
    }
    //defer client.Close()

    n, err := client.Write([]byte("list"))
    if err != nil {
        panic(err)
    }

    buffer := make([]byte, 512)
    n, err2 := client.Read(buffer)

    var slice []string
    json.Unmarshal(buffer[:n], slice)
    if err2 != nil {
        fmt.Println("Read failed:", err2)
        return
    }
    fmt.Println( "msg:", string(buffer[:n]))
    err3 := client.Close()
    if err!=nil{
        fmt.Println(err3)
    }
}