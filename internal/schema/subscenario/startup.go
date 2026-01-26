package subscenarioschema

type StartUpReq struct {
	Module        string `json:"module"`
	Action        string `json:"action"`
	TimeStamp     int    `json:"timeStamp"`
	SubscenarioID int    `json:"subscenario_id"`
	Mgd           int    `json:"mgd"`
	CommandNum    string `json:"commandNum"`
}

type StartUpData struct {
	SubscenarioID      int   `json:"subscenario_id"`
	ScenarioAdjustment int   `json:"scenario_adjustment"`
	ServerTimestamp    int64 `json:"server_timestamp"`
}

type StartUpResp struct {
	ResponseData StartUpData `json:"response_data"`
	ReleaseInfo  []any       `json:"release_info"`
	StatusCode   int         `json:"status_code"`
}
