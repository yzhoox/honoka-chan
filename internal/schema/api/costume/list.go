package costumeapischema

type CostumeList struct {
	UnitID    int  `json:"unit_id"`
	IsRankMax bool `json:"is_rank_max"`
	IsSigned  bool `json:"is_signed"`
}

type ListData struct {
	CostumeList []CostumeList `json:"costume_list"`
}

type ListResp struct {
	Result     ListData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
