package ali_drive

type AliRefreshRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type ApiRefreshResponse struct {
	AccessToken    string   `json:"access_token"`
	Avatar         string   `json:"avatar"`
	DefaultDriveId string   `json:"default_drive_id"`
	DeviceId       string   `json:"device_id"`
	ExistLink      []string `json:"exist_link"`
	ExpireTime     string   `json:"expire_time"`
	ExpiresIn      int      `json:"expires_in"`
	IsFirstLogin   bool     `json:"is_first_login"`
	NeedLink       bool     `json:"need_link"`
	NeedRpVerify   bool     `json:"need_rp_verify"`
	NickName       string   `json:"nick_name"`
	PinSetup       bool     `json:"pin_setup"`
	RefreshToken   string   `json:"refresh_token"`
	Role           string   `json:"role"`
	State          string   `json:"state"`
	Status         string   `json:"status"`
	TokenType      string   `json:"token_type"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
}

func (y *Ali) Refresh(refreshToken string) *ApiRefreshResponse {
	url := "https://auth.aliyundrive.com/v2/account/token"
	req := y.r().
		SetBody(AliRefreshRequest{
			GrantType:    "refresh_token",
			RefreshToken: refreshToken,
		})

	if !y.unmarshal(y.processResp(req.Post(url)), &y.RefreshInfo) {
		return nil
	}
	return &y.RefreshInfo
}
