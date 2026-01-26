package unitschema

type SetDisplayRankReq struct {
	Module           string `json:"module"`
	Action           string `json:"action"`
	TimeStamp        int64  `json:"timeStamp"`
	Mgd              int    `json:"mgd"`
	UnitOwningUserID int    `json:"unit_owning_user_id"`
	DisplayRank      int    `json:"display_rank"`
	CommandNum       string `json:"commandNum"`
}

type SetDisplayRankResp struct {
	ResponseData []any `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}
