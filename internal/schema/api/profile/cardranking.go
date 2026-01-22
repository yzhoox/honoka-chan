package profile

type CardRankingResp struct {
	Result     []any `json:"result"`
	Status     int   `json:"status"`
	CommandNum bool  `json:"commandNum"`
	TimeStamp  int64 `json:"timeStamp"`
}
