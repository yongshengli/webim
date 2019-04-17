package server

import (
    "github.com/astaxie/beego/logs"
    "webim/comet/common"
    "webim/comet/models"
)

func (j *JobWorker) leaveRoom() {
    var reqData map[string]interface{}
    if room := j.getRoom(&reqData); room != nil {
        room.Leave(j.s)

        j.Rsp.Type = TYPE_ROOM_MSG
        rspData := map[string]interface{}{"code": 0, "content": "ok"}
        resByte, err := common.EnJson(rspData)
        if err != nil {
            logs.Error("msg[leaveRoom encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    }
}
func (j *JobWorker) joinRoom() {
    var reqData map[string]interface{}
    room := j.getRoom(&reqData)
    if room == nil {
        rspData := make(map[string]interface{})
        rspData["code"] = 1
        rspData["content"] = "房间不存在"
        rspData["room_id"] = reqData["room_id"]

        resByte, err := common.EnJson(rspData)
        if err != nil {
            logs.Error("msg[joinRoom encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Type = TYPE_ROOM_MSG
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    } else {
        _, err := room.Join(j.s)
        if err != nil {
            logs.Error("msg[%s]", err.Error())
            return
        }

        rspData := make(map[string]interface{})
        rspData["code"] = 1
        rspData["content"] = j.s.User.Name + "进入房间"
        resByte, err := common.EnJson(rspData)
        if err != nil {
            logs.Error("msg[joinRoom encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Data = string(resByte)
        room.Broadcast(&j.Rsp)

        j.sendLastChatToCurrentUser(room)

    }
}
func (j *JobWorker) roomMsg() {
    var rspData map[string]interface{}
    room := j.getRoom(&rspData)
    if room == nil {
        rspData := make(map[string]interface{})
        rspData["content"] = "房间不存在"
        rspData["room_id"] = rspData["room_id"]

        j.Rsp.Type = TYPE_ROOM_MSG
        resByte, err := common.EnJson(rspData)
        if err != nil {
            logs.Error("msg[roomMsg encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    } else {
        rspData["uid"] = j.s.User.Id
        rspData["room_id"] = room.Id
        if TmpRspData, err := common.EnJson(rspData); err==nil {
            j.Rsp.Data = string(TmpRspData)
            room.Broadcast(&j.Rsp)
        }else{
            logs.Error("msg[roomMsg EnJson err] err[%s]", err)
            return
        }
    }
}

func (j *JobWorker) createRoom() {
    var reqData map[string]interface{}
    room := j.getRoom(&reqData)
    if room == nil {
        room, _ = NewRoom(reqData["room_id"].(string), "")
    }
    _, err := room.Join(j.s)
    if err != nil {
        logs.Error("msg[加入房间失败] err[%s]", err.Error())
        return
    }
    rspData := make(map[string]interface{})
    rspData["code"] = 0
    rspData["content"] = "创建房间成功"
    j.Rsp.Type = TYPE_CREATE_ROOM
    resByte, err := common.EnJson(rspData)
    if err != nil {
        logs.Error("msg[roomMsg encode err] err[%s]", err.Error())
        return
    }
    j.Rsp.Data = string(resByte)
    j.s.Send(&j.Rsp)

    //发送历史聊天记录
    j.sendLastChatToCurrentUser(room)
}

func (j *JobWorker) sendLastChatToCurrentUser(room *Room) error {
    msgArr, err := models.FindRoomMsgLast(room.Id, 30)
    if err != nil {
        logs.Error("msg[获取聊天室最后30条聊天记录失败] err[%s]", err.Error())
        return err
    }
    sortMsgArr := make([]models.RoomMsg, len(msgArr))
    jj := 0
    for i:=len(msgArr)-1; i>0; i++ {
        sortMsgArr[jj] = msgArr[i]
        jj++
    }
    rspData := map[string]interface{}{}
    rspData["room_id"] = room.Id
    rspData["chat_history"] = sortMsgArr
    msgArr = nil
    tmpByte, err := common.EnJson(rspData)
    if err != nil {
        logs.Error("msg[sendLastChatToCurrentUser json encode err] err[%s]", err.Error())
        return err
    }
    j.Rsp.Data = string(tmpByte)
    return nil
}

func (j *JobWorker) getRoom(reqData *map[string]interface{}) *Room{
    var err error
    *reqData, err = j.decode(j.Req.Data)
    if err != nil {
        logs.Error("msg[leaveRoom decode err] err[%s]", err.Error())
        return nil
    }
    if _, ok := (*reqData)["room_id"]; !ok {
        logs.Warn("msg[room_id为空]")
        return nil
    }
    room, err := GetRoom((*reqData)["room_id"].(string))
    if err != nil {
        logs.Error("msg[获取房价失败] err[%s]", err.Error())
        return nil
    }
    return room
}