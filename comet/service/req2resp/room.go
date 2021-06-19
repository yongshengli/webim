package req2resp

import "comet/models"

type LeaveRoomReq struct {
	RoomId string `json:"room_id"`
}

type LeaveRoomResp struct {
	Status Status `json:"status"`
}

type JoinRoomReq struct {
	RoomId string `json:"room_id"`
}
type JoinRoomResp struct {
	Status Status `json:"status"`
}

type CreateRoomReq struct {
	RoomId string `json:"room_id"`
}

type CreateRoomResp struct {
	Status Status `json:"status"`
}

type RoomMsgReq struct {
	RoomId  string `json:"room_id"`
	Content string `json:"content"`
}

type RoomMsgResp struct {
	Status Status `json:"status"`
	Body   models.RoomMsg
}
