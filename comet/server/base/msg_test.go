package base

import (
	"testing"
)

func TestMap2Msg(t *testing.T) {
	arr := []map[string]interface{}{
		{"type": float64(10), "data": ""},
		{"type": float64(1), "data": ""},
	}
	for _, m := range arr {
		msg, err := Map2Msg(m)
		if err != nil {
			t.Error(err)
		}
		//fmt.Println(msg)
		if msg.Type == 0 {
			t.Errorf("转换后的type是%f,期望是%f", msg.Type, msg.Type)
		}
	}
}
