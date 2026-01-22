package login

type LoginData struct {
	AuthorizeToken  string `json:"authorize_token"`
	UserId          int    `json:"user_id"`
	ReviewVersion   string `json:"review_version"`
	ServerTimestamp int64  `json:"server_timestamp"`
	IdfaEnabled     bool   `json:"idfa_enabled"`
	SkipLoginNews   bool   `json:"skip_login_news"`
	AdultFlag       int    `json:"adult_flag"`
}

type LoginResp struct {
	ResponseData LoginData `json:"response_data"`
	ReleaseInfo  []any     `json:"release_info"`
	StatusCode   int       `json:"status_code"`
}
