package server

import (
    "github.com/astaxie/beego/logs"
    "webim/comet/common"
)

func (j *JobWorker) leaveRoom() {
    var data map[string]interface{}
    if room := j.getRoom(&data); room != nil {
        room.Leave(j.s)

        j.Rsp.Type = TYPE_ROOM_MSG
        resData := map[string]interface{}{"code": 0, "content": "ok"}
        resByte, err := common.EnJson(resData)
        if err != nil {
            logs.Error("msg[leaveRoom encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    }
}
func (j *JobWorker) joinRoom() {
    var data map[string]interface{}
    room := j.getRoom(&data)
    if room == nil {
        resData := make(map[string]interface{})
        resData["code"] = 1
        resData["content"] = "房间不存在"
        resData["room_id"] = data["room_id"]

        resByte, err := common.EnJson(resData)
        if err != nil {
            logs.Error("msg[joinRoom encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Type = TYPE_ROOM_MSG
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    } else {
        res, err := room.Join(j.s)
        if err != nil {
            logs.Error("msg[%s]", err.Error())
            return
        }
        if res {
            resData := make(map[string]interface{})
            resData["code"] = 1
            resData["content"] = j.s.User.Name + "进入房间"
            resByte, err := common.EnJson(resData)
            if err != nil {
                logs.Error("msg[joinRoom encode err] err[%s]", err.Error())
                return
            }
            j.Rsp.Data = string(resByte)
            room.Broadcast(&j.Rsp)
        }
    }
}
func (j *JobWorker) roomMsg() {
    var data map[string]interface{}
    room := j.getRoom(&data)
    if room == nil {
        resData := make(map[string]interface{})
        resData["content"] = "房间不存在"
        resData["room_id"] = data["room_id"]

        j.Rsp.Type = TYPE_ROOM_MSG
        resByte, err := common.EnJson(resData)
        if err != nil {
            logs.Error("msg[roomMsg encode err] err[%s]", err.Error())
            return
        }
        j.Rsp.Data = string(resByte)
        j.s.Send(&j.Rsp)
    } else {
        data["uid"] = j.s.User.Id
        data["room_id"] = room.Id
        if TmpRspData, err := common.EnJson(data); err==nil {
            j.Rsp.Data = string(TmpRspData)
            room.Broadcast(&j.Rsp)
        }else{
            logs.Error("msg[roomMsg EnJson err] err[%s]", err)
            return
        }
    }
}
func (j *JobWorker) createRoom() {
    var data map[string]interface{}
    room := j.getRoom(&data)
    if room == nil {
        room, _ = NewRoom(data["room_id"].(string), "")
    }
    _, err := room.Join(j.s)
    if err != nil {
        logs.Error("msg[加入房间失败] err[%s]", err.Error())
        return
    }
    resData := make(map[string]interface{})
    resData["code"] = 0
    resData["content"] = "创建房间成功"
    j.Rsp.Type = TYPE_CREATE_ROOM
    resByte, err := common.EnJson(resData)
    if err != nil {
        logs.Error("msg[roomMsg encode err] err[%s]", err.Error())
        return
    }
    j.Rsp.Data = string(resByte)
    j.s.Send(&j.Rsp)
}

func (j *JobWorker) getRoom(data *map[string]interface{}) *Room{
    var err error
    *data, err = j.decode(j.Req.Data)
    if err != nil {
        logs.Error("msg[leaveRoom decode err] err[%s]", err.Error())
        return nil
    }
    if _, ok := (*data)["room_id"]; !ok {
        logs.Warn("msg[room_id为空]")
        return nil
    }
    room, err := GetRoom((*data)["room_id"].(string))
    if err != nil {
        logs.Error("msg[获取房价失败] err[%s]", err.Error())
        return nil
    }
    return room
}