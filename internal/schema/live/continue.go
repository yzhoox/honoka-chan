package liveschema

type ContinueData struct {
	BeforeSnsCoin   int   `json:"before_sns_coin"`
	AfterSnsCoin    int   `json:"after_sns_coin"`
	ServerTimestamp int64 `json:"server_timestamp"`
}

type ContinueResp struct {
	ResponseData ContinueData `json:"response_data"`
	ReleaseInfo  []any        `json:"release_info"`
	StatusCode   int          `json:"status_code"`
}
