package personalnoticeschema

type GetData struct {
	HasNotice       bool   `json:"has_notice"`
	NoticeID        int    `json:"notice_id"`
	Type            int    `json:"type"`
	Title           string `json:"title"`
	Contents        string `json:"contents"`
	ServerTimestamp int64  `json:"server_timestamp"`
}

type GetResp struct {
	ResponseData GetData `json:"response_data"`
	ReleaseInfo  []any   `json:"release_info"`
	StatusCode   int     `json:"status_code"`
}
