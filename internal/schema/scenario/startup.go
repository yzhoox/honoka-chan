package scenario

type StartUpReq struct {
	Module     string `json:"module"`
	Action     string `json:"action"`
	TimeStamp  int    `json:"timeStamp"`
	Mgd        int    `json:"mgd"`
	CommandNum string `json:"commandNum"`
	ScenarioID int    `json:"scenario_id"`
}

type StartUpData struct {
	ScenarioID         int   `json:"scenario_id"`
	ScenarioAdjustment int   `json:"scenario_adjustment"`
	ServerTimestamp    int64 `json:"server_timestamp"`
}

type StartUpResp struct {
	ResponseData StartUpData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
