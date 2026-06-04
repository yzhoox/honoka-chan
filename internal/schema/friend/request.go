package friendschema

type RequestData struct {
	IsFriend bool `json:"is_friend"`
}

type RequestResp struct {
	ResponseData RequestData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
