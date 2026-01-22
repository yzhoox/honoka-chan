package navigation

type SpecialCutinData struct {
	SpecialCutinList []any `json:"special_cutin_list"`
}

type SpecialCutinResp struct {
	Result     SpecialCutinData `json:"result"`
	Status     int              `json:"status"`
	CommandNum bool             `json:"commandNum"`
	TimeStamp  int64            `json:"timeStamp"`
}
