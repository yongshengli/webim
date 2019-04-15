package models

type RoomMsg struct {
    Id      int    `orm:"auto"`
    RoomId  int
    Uid     int
    Content string `orm:"size(100)"`
    CT      int    `orm:"auto_now_add;type(int32)"`
}

func (rm * RoomMsg) TableName() string {
    return "room_msg_0"
}