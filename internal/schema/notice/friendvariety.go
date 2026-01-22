package notice

type FriendVarietyData struct {
	ItemCount       int   `json:"item_count"`
	NoticeList      []any `json:"notice_list"`
	ServerTimestamp int64 `json:"server_timestamp"`
}

type FriendVarietyResp struct {
	ResponseData FriendVarietyData `json:"response_data"`
	ReleaseInfo  []any             `json:"release_info"`
	StatusCode   int               `json:"status_code"`
}
