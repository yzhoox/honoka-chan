package backgroundapischema

type Info struct {
	BackgroundID int    `json:"background_id"`
	IsSet        bool   `json:"is_set"`
	InsertDate   string `json:"insert_date"`
}

type InfoData struct {
	BackgroundInfo []Info `json:"background_info"`
}

type InfoResp struct {
	Result     InfoData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
