package friendschema

type RequestCancelData struct {
	IsFriend bool `json:"is_friend"`
}

type RequestCancelResp struct {
	ResponseData RequestCancelData `json:"response_data"`
	ReleaseInfo  []any             `json:"release_info"`
	StatusCode   int               `json:"status_code"`
}
