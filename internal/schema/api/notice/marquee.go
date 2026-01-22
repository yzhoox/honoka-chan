package notice

type MarqueeData struct {
	ItemCount   int   `json:"item_count"`
	MarqueeList []any `json:"marquee_list"`
}

type MarqueeResp struct {
	Result     MarqueeData `json:"result"`
	Status     int         `json:"status"`
	CommandNum bool        `json:"commandNum"`
	TimeStamp  int64       `json:"timeStamp"`
}
