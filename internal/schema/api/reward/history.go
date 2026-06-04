package rewardapischema

type HistoryData struct {
	ItemCount int   `json:"item_count"`
	Limit     int   `json:"limit"`
	History   []any `json:"history"`
}

type HistoryResp struct {
	Result     HistoryData `json:"result"`
	Status     int         `json:"status"`
	CommandNum bool        `json:"commandNum"`
	TimeStamp  int64       `json:"timeStamp"`
}
