package award

type Info struct {
	AwardID    int    `json:"award_id"`
	IsSet      bool   `json:"is_set"`
	InsertDate string `json:"insert_date"`
}

type InfoData struct {
	AwardInfo []Info `json:"award_info"`
}

type InfoResp struct {
	Result     InfoData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
