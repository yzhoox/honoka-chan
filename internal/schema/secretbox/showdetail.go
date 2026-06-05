package secretboxschema

type ShowDetailReq struct {
	Module      string `json:"module"`
	Action      string `json:"action"`
	TimeStamp   int64  `json:"timeStamp"`
	Mgd         int    `json:"mgd"`
	SecretBoxID int    `json:"secret_box_id"`
	CommandNum  string `json:"commandNum"`
}

type UnitLineUp struct {
	Rarity  int   `json:"rarity"`
	UnitIDs []int `json:"unit_ids"`
}

type ButtonTypeUnitLineUp struct {
	SecretBoxButtonType int          `json:"secret_box_button_type"`
	SecretBoxName       string       `json:"secret_box_name"`
	UnitLineUp          []UnitLineUp `json:"unit_line_up"`
}

type ShowDetailData struct {
	ButtonTypeUnitLineUp []ButtonTypeUnitLineUp `json:"button_type_unit_line_up,omitempty"`
	UnitLineUp           []UnitLineUp           `json:"unit_line_up,omitempty"`
	URL                  string                 `json:"url"`
}

type ShowDetailResp struct {
	ResponseData ShowDetailData `json:"response_data"`
	ReleaseInfo  []any          `json:"release_info"`
	StatusCode   int            `json:"status_code"`
}
