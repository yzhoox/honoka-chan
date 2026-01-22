package notice

type FriendGreetingData struct {
	NextId          int   `json:"next_id"`
	NoticeList      []any `json:"notice_list"`
	ServerTimestamp int64 `json:"server_timestamp"`
}

type FriendGreetingResp struct {
	ResponseData FriendGreetingData `json:"response_data"`
	ReleaseInfo  []any              `json:"release_info"`
	StatusCode   int                `json:"status_code"`
}
