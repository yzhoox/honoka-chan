package eventscenarioschema

type StartUpReq struct {
	Module          string `json:"module"`
	EventScenarioID int    `json:"event_scenario_id"`
	Action          string `json:"action"`
	TimeStamp       int    `json:"timeStamp"`
	Mgd             int    `json:"mgd"`
	CommandNum      string `json:"commandNum"`
}

type EventScenarioList struct {
	EventID         int `json:"event_id"`
	Progress        int `json:"progress"`
	Status          int `json:"status"`
	EventScenarioID int `json:"event_scenario_id"`
}

type StartUpData struct {
	EventScenarioList  EventScenarioList `json:"event_scenario_list"`
	ScenarioAdjustment int               `json:"scenario_adjustment"`
}

type StartUpResp struct {
	ResponseData StartUpData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
