package announce

type CheckStateData struct {
	HasUnreadAnnounce bool  `json:"has_unread_announce"`
	ServerTimestamp   int64 `json:"server_timestamp"`
}

type CheckStateDataResp struct {
	ResponseData CheckStateData `json:"response_data"`
	ReleaseInfo  []any          `json:"release_info"`
	StatusCode   int            `json:"status_code"`
}
