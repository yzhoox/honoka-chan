package eventscenarioapischema

type ChapterList struct {
	EventScenarioID int    `json:"event_scenario_id"`
	Chapter         int    `json:"chapter"`
	ChapterAsset    string `json:"chapter_asset,omitempty"`
	Status          int    `json:"status"`
	OpenFlashFlag   int    `json:"open_flash_flag"`
	IsReward        bool   `json:"is_reward"`
	CostType        int    `json:"cost_type"`
	ItemID          int    `json:"item_id"`
	Amount          int    `json:"amount"`
}

type EventScenarioList struct {
	EventID               int           `json:"event_id"`
	EventScenarioBtnAsset string        `json:"event_scenario_btn_asset"`
	OpenDate              string        `json:"open_date"`
	ChapterList           []ChapterList `json:"chapter_list"`
}

type StatusData struct {
	EventScenarioList []EventScenarioList `json:"event_scenario_list"`
}

type StatusResp struct {
	Result     StatusData `json:"result"`
	Status     int        `json:"status"`
	CommandNum bool       `json:"commandNum"`
	TimeStamp  int64      `json:"timeStamp"`
}
