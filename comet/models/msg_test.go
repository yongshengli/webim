package models

import (
    "testing"
    "reflect"
    "fmt"
)

func TestMapToMsg(t *testing.T) {
    m := map[string]interface{}{"type":float64(10), "data":map[string]interface{}{}}

    fmt.Println(reflect.ValueOf(float64(10)).Kind())
    msg := Map2Msg(m)
    //fmt.Println(msg)
    if m["type"]!=msg.Type {
        t.Errorf("转换后的type是%d,期望是%d", msg.Type, m["type"])
    }
}
