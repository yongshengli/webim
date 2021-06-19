package req2resp

type Status struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
}

type LoginResp struct {
	Status Status `json:"status"`
	Token  string `json:"token"`
}
