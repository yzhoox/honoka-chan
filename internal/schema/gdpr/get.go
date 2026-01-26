package gdprschema

type GetData struct {
	EnableGdpr      bool  `json:"enable_gdpr"`
	IsEea           bool  `json:"is_eea"`
	ServerTimestamp int64 `json:"server_timestamp"`
}

type GetResp struct {
	ResponseData GetData `json:"response_data"`
	ReleaseInfo  []any   `json:"release_info"`
	StatusCode   int     `json:"status_code"`
}
