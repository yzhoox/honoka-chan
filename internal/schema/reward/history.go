package rewardschema

type HistoryData struct {
	ItemCount int   `json:"item_count"`
	Limit     int   `json:"limit"`
	History   []any `json:"history"`
}

type HistoryResp struct {
	ResponseData HistoryData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
