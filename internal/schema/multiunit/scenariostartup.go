package multiunitschema

type ScenarioStartUpReq struct {
	Module              string `json:"module"`
	Action              string `json:"action"`
	TimeStamp           int    `json:"timeStamp"`
	Mgd                 int    `json:"mgd"`
	MultiUnitScenarioID int    `json:"multi_unit_scenario_id"`
	CommandNum          string `json:"commandNum"`
}

type ScenarioStartUpData struct {
	MultiUnitScenarioID int   `json:"multi_unit_scenario_id"`
	ScenarioAdjustment  int   `json:"scenario_adjustment"`
	ServerTimestamp     int64 `json:"server_timestamp"`
}

type ScenarioStartUpResp struct {
	ResponseData ScenarioStartUpData `json:"response_data"`
	ReleaseInfo  []any               `json:"release_info"`
	StatusCode   int                 `json:"status_code"`
}
