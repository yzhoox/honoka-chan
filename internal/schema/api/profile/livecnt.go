package profileapischema

type LiveCntData struct {
	Difficulty int `json:"difficulty"`
	ClearCnt   int `json:"clear_cnt"`
}

type LiveCntResp struct {
	Result     []LiveCntData `json:"result"`
	Status     int           `json:"status"`
	CommandNum bool          `json:"commandNum"`
	TimeStamp  int64         `json:"timeStamp"`
}
