package unit

type UnitDeckDetail struct {
	Position         int `json:"position"`
	UnitOwningUserID int `json:"unit_owning_user_id"`
}

type UnitDeckList struct {
	UnitDeckDetail []UnitDeckDetail `json:"unit_deck_detail"`
	UnitDeckID     int              `json:"unit_deck_id"`
	MainFlag       int              `json:"main_flag"`
	DeckName       string           `json:"deck_name"`
}

type DeckReq struct {
	Module       string         `json:"module"`
	UnitDeckList []UnitDeckList `json:"unit_deck_list"`
	Action       string         `json:"action"`
	Mgd          int            `json:"mgd"`
}

type DeckResp struct {
	ResponseData []any `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}
