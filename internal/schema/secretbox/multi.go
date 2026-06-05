package secretboxschema

type MultiReq struct {
	Module      string `json:"module"`
	Action      string `json:"action"`
	Mgd         int    `json:"mgd"`
	ID          int    `json:"id"`
	SecretBoxID int    `json:"secret_box_id"`
	Count       int    `json:"count"`
	TimeStamp   int64  `json:"timeStamp"`
	CommandNum  string `json:"commandNum"`
	UnitTypeIds []any  `json:"unit_type_ids"`
}

type MultiResp struct {
	ResponseData PonData `json:"response_data"`
	ReleaseInfo  []any   `json:"release_info"`
	StatusCode   int     `json:"status_code"`
}
