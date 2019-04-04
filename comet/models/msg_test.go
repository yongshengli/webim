package models

import (
    "testing"
)

func TestMapToMsg(t *testing.T) {
    m := map[string]interface{}{"type":TYPE_NOTICE_MSG, "data":map[string]interface{}{}}

    msg := Map2Msg(m)
    if m["type"]!=msg.Type {
        t.Errorf("转换后的type是%d,期望是%d", msg.Type, m["type"])
    }
}
