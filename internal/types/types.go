package types

const AppVersion = "1.5.9"

// UnlockData represents stored unlock data
type UnlockData struct {
	User     string `json:"user"`
	Password string `json:"pwd"`
	WbID     string `json:"wb_id"`
	Login    string `json:"login"`
	UID      string `json:"uid"`
}

// DeviceInfo represents device information
type DeviceInfo struct {
	Unlocked string
	Product  string
	SoC      string
	Token    string
}

// XiaomiAuthResponse represents Xiaomi authentication response
type XiaomiAuthResponse struct {
	Code            int    `json:"code"`
	SecurityStatus  int    `json:"securityStatus"`
	NotificationURL string `json:"notificationUrl"`
	SSecurity       string `json:"ssecurity"`
	Nonce           string `json:"nonce"`
	Location        string `json:"location"`
	PassToken       string `json:"passToken"`
	UserID          string `json:"userId"`
}

// UnlockResponse represents unlock API response
type UnlockResponse struct {
	Code        int    `json:"code"`
	DescEN      string `json:"descEN"`
	EncryptData string `json:"encryptData"`
	Data        struct {
		WaitHour int `json:"waitHour"`
	} `json:"data"`
}
