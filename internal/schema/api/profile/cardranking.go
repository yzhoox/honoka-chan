package profileapischema

type CardRankingData struct {
	UnitID    int  `json:"unit_id"`
	TotalLove int  `json:"total_love"`
	Rank      int  `json:"rank"`
	SignFlag  bool `json:"sign_flag"`
}

type CardRankingResp struct {
	Result     any   `json:"result"`
	Status     int   `json:"status"`
	CommandNum bool  `json:"commandNum"`
	TimeStamp  int64 `json:"timeStamp"`
}
