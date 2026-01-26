package noticeschema

type UserGreetingData struct {
	ItemCount       int   `json:"item_count"`
	HasNext         bool  `json:"has_next"`
	NoticeList      []any `json:"notice_list"`
	ServerTimestamp int64 `json:"server_timestamp"`
}

type UserGreetingResp struct {
	ResponseData UserGreetingData `json:"response_data"`
	ReleaseInfo  []any            `json:"release_info"`
	StatusCode   int              `json:"status_code"`
}
