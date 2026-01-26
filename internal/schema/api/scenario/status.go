package scenarioapischema

type StatusList struct {
	ScenarioID int `json:"scenario_id"`
	Status     int `json:"status"`
}

type StatusData struct {
	ScenarioStatusList []StatusList `json:"scenario_status_list"`
}

type StatusResp struct {
	Result     StatusData `json:"result"`
	Status     int        `json:"status"`
	CommandNum bool       `json:"commandNum"`
	TimeStamp  int64      `json:"timeStamp"`
}
