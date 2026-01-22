package subscenario

type StatusList struct {
	SubscenarioID int `json:"subscenario_id"`
	Status        int `json:"status"`
}

type StatusData struct {
	SubscenarioStatusList  []StatusList `json:"subscenario_status_list"`
	UnlockedSubscenarioIds []any        `json:"unlocked_subscenario_ids"`
}

type StatusResp struct {
	Result     StatusData `json:"result"`
	Status     int        `json:"status"`
	CommandNum bool       `json:"commandNum"`
	TimeStamp  int64      `json:"timeStamp"`
}
