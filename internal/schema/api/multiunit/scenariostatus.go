package multiunitapischema

type ChapterList struct {
	MultiUnitScenarioID int `json:"multi_unit_scenario_id"`
	Chapter             int `json:"chapter"`
	Status              int `json:"status"`
}

type StatusList struct {
	MultiUnitID               int           `json:"multi_unit_id"`
	Status                    int           `json:"status"`
	MultiUnitScenarioBtnAsset string        `json:"multi_unit_scenario_btn_asset"`
	OpenDate                  string        `json:"open_date"`
	ChapterList               []ChapterList `json:"chapter_list"`
}

type StatusData struct {
	MultiUnitScenarioStatusList  []StatusList `json:"multi_unit_scenario_status_list"`
	UnlockedMultiUnitScenarioIds []any        `json:"unlocked_multi_unit_scenario_ids"`
}

type StatusResp struct {
	Result     StatusData `json:"result"`
	Status     int        `json:"status"`
	CommandNum bool       `json:"commandNum"`
	TimeStamp  int64      `json:"timeStamp"`
}
