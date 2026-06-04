package rewardapischema

type ListData struct {
	ItemCount int   `json:"item_count"`
	Limit     int   `json:"limit"`
	Order     int   `json:"order"`
	Items     []any `json:"items"`
}

type ListResp struct {
	Result     ListData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
