package livese

type InfoData struct {
	LiveSeList []int `json:"live_se_list"`
}

type InfoResp struct {
	Result     InfoData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
