package actions

import (
	"comet/common"
	"comet/models"
	"comet/server/base"
	"comet/server/req2resp"
	"comet/server/room"

	"github.com/astaxie/beego/logs"
)

type Action struct {
}

func (j *Action) LeaveRoom(req, resp *base.Msg) error {
	var reqData req2resp.LeaveRoomReq
	common.DeJson([]byte(req.Data), &reqData)
	if tRoom := j.getRoom(reqData.RoomId); tRoom != nil {
		tRoom.Leave(j.s)
		rspData := &req2resp.LeaveRoomResp{
			Status: req2resp.Status{
				Code: 0,
				Msg:  "ok",
			},
		}
		data, _ := common.EnJson(rspData)
		resp.Data = string(data)
	}
	return nil
}
func (j *Action) JoinRoom(req, resp *base.Msg) error {
	var reqData req2resp.JoinRoomReq
	common.DeJson([]byte(req.Data), &reqData)
	tRoom := j.getRoom(reqData.RoomId)
	if tRoom == nil {
		rspData := &req2resp.JoinRoomResp{
			Status: req2resp.Status{
				Code: 1,
				Msg:  "房间不存在",
			},
		}
		tmpResp, _ := common.EnJson(rspData)
		resp.Data = string(tmpResp)
	} else {
		_, err := tRoom.Join(j.s)
		if err != nil {
			logs.Error("msg[%s]", err.Error())
			return nil
		}
		roomMsg := &models.RoomMsg{
			RoomId:  reqData.RoomId,
			Content: j.s.User.Name + "进入房间",
		}
		room.Broadcast(roomMsg)
		rspData := &req2resp.JoinRoomResp{
			Status: req2resp.Status{
				Code: 1,
			},
		}
		tmpData, _ := common.EnJson(rspData)
		resp.Data = string(tmpData)
		j.sendLastChatToCurrentUser(tRoom)
	}
	return nil
}
func (j *Action) RoomMsg(req, resp *base.Msg) error {
	var reqData req2resp.RoomMsgReq
	common.DeJson([]byte(req.Data), &reqData)
	tRoom := j.getRoom(reqData.RoomId)
	rspData := &req2resp.RoomMsgResp{}
	if tRoom == nil {
		rspData.Status = req2resp.Status{
			Code: 1,
			Msg:  "房间不存在",
		}
	} else {
		rspData.Body.Uid = j.s.User.Id
		rspData.Body.Uname = j.s.User.Name
		rspData.Body.RoomId = room.Id
		if TmpRspData, err := common.EnJson(rspData); err == nil {
			j.Rsp.Data = string(TmpRspData)
			room.Broadcast(&rspData.Body)
		} else {
			logs.Error("msg[roomMsg EnJson err] err[%s]", err)
			return nil
		}
	}
	tmpResp, _ := common.EnJson(rspData)
	resp.Data = string(tmpResp)
	return nil
}

func (j *Action) CreateRoom(req, resp *base.Msg) {
	var reqData req2resp.CreateRoomReq
	common.DeJson([]byte(req.Data), &reqData)
	tRoom := j.getRoom(reqData.RoomId)
	if tRoom == nil {
		tRoom, _ = room.NewRoom(reqData.RoomId, "")
	}
	_, err := tRoom.Join(j.s)
	if err != nil {
		logs.Error("msg[加入房间失败] err[%s]", err.Error())
		return
	}
	rspData := req2resp.CreateRoomResp{
		Status: req2resp.Status{
			Code: 0,
			Msg:  "创建房间成功",
		},
	}
	j.Rsp.Type = base.TYPE_CREATE_ROOM
	resByte, err := common.EnJson(rspData)
	if err != nil {
		logs.Error("msg[roomMsg encode err] err[%s]", err.Error())
		return
	}
	j.Rsp.Data = string(resByte)
	//j.s.Send(j.Rsp)

	//发送历史聊天记录
	j.sendLastChatToCurrentUser(tRoom)
}

func (j *Action) sendLastChatToCurrentUser(room *room.Room) error {
	msgArr, err := models.GetLastRoomMsg(room.Id, 30)
	if err != nil {
		logs.Error("msg[获取聊天室最后30条聊天记录失败] err[%s]", err.Error())
		return err
	}
	sortMsgArr := make([]models.RoomMsg, len(msgArr))
	jj := 0
	for i := len(msgArr) - 1; i >= 0; i-- {
		sortMsgArr[jj] = msgArr[i]
		jj++
	}
	rspData := map[string]interface{}{}

	if err = common.DeJson([]byte(j.Rsp.Data), &rspData); err != nil {
		logs.Error("msg[sendLastChatToCurrentUser json decode err] err[%s]", err.Error())
		return err
	}
	msgArr = nil
	rspData["chat_history"] = sortMsgArr
	j.Rsp.Data, _ = common.Map2String(rspData)
	j.s.Send(j.Rsp)
	return nil
}

func (j *Action) getRoom(roomId string) *room.Room {
	var err error
	tRoom, err := room.GetRoom(roomId)
	if err != nil {
		logs.Error("msg[获取房间失败] err[%s]", err.Error())
		return nil
	}
	return tRoom
}
