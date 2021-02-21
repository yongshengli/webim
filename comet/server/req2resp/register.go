package req2resp

type RegisterReq struct {
	DeviceId string `json:"device_id"`
}
type RegisterResp struct {
	Status   Status `json:"status"`
	DeviceId string `json:"device_id"`
}
