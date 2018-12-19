package models

const (
	MSG_ROOM = iota
	MSG_COMMON
)

type Msg struct {
	MsgType   int                         `json:"type"`
	Data      map[interface{}]interface{} `json:"data"`
	Timestamp int                         `json:"time"`
}