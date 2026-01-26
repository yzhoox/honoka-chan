package unitschema

type FavoriteReq struct {
	Module           string `json:"module"`
	Action           string `json:"action"`
	TimeStamp        int64  `json:"timeStamp"`
	Mgd              int    `json:"mgd"`
	FavoriteFlag     int    `json:"favorite_flag"`
	UnitOwningUserID int    `json:"unit_owning_user_id"`
	CommandNum       string `json:"commandNum"`
}

type FavoriteResp struct {
	ResponseData []any `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}
