package userschema

type ChangeNameData struct {
	BeforeName      string `json:"before_name"`
	AfterName       string `json:"after_name"`
	ServerTimestamp int64  `json:"server_timestamp"`
}

type ChangeNameResp struct {
	ResponseData ChangeNameData `json:"response_data"`
	ReleaseInfo  []any          `json:"release_info"`
	StatusCode   int            `json:"status_code"`
}
