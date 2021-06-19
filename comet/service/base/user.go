package base

type User struct {
	Id            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Platform      string                 `json:"platform"`
	ClientVersion string                 `json:"clientVersion"`
	DeviceId      string                 `json:"device_id"`
	DeviceToken   string                 `json:"device_token"` // CometToken = md5(udid+appKey)
	Info          map[string]interface{} `json:"info"`
	IP            string                 `json:"ip"`
	RealIP        string                 `json:"real_ip"`
}
