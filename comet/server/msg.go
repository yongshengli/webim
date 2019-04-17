package server

import (
	"reflect"
	"errors"
	"fmt"
)


type Job struct {
	TraceID string // 跟踪ID

	ReqTime int64 // 请求到达时间
	RspTime int64 // 响应结束时间

	TransferReqTime int64 // 转发请求开始时间
	TransferRspTime int64 // 转发请求结束时间

	Req Msg // 原请求信息
	Rsp Msg // 响应信息
}

type Msg struct {
	Type        float64 `json:"type" valid:"Required"`
	DeviceToken string  `json:"device_token"`
	Version     string  `json:"version"`
	ReqId       string  `json:"req_id"`
	Encode      string  `json:"encode"`
	Data        string  `json:"data"`
}

func Map2Msg(data map[string]interface{}) (Msg, error){
	msg := &Msg{}
	elem := reflect.ValueOf(msg).Elem()
	relType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		tag := relType.Field(i).Tag
		mk := tag.Get("json")
		if mv, ok := data[mk]; ok {
			if elem.Field(i).Kind() != reflect.ValueOf(mv).Kind() {
				return *msg, errors.New(fmt.Sprintf("Map2Msg error: map[%s]type error", mk))
			}
			elem.Field(i).Set(reflect.ValueOf(mv))
		}
	}
	return *msg, nil
}