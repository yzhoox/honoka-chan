package apischema

type ApiReq struct {
	Module    string `json:"module"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timeStamp"`
}

type ApiResp struct {
	ResponseData any `json:"response_data"`
	ReleaseInfo  any `json:"release_info"`
	StatusCode   int `json:"status_code"`
}
