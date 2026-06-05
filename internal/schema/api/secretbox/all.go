package secretboxapischema

type ItemAmount struct {
	ItemID int `json:"item_id"`
	Amount int `json:"amount"`
}

type GaugeInfo struct {
	MaxGaugePoint int `json:"max_gauge_point"`
	GaugePoint    int `json:"gauge_point"`
}

type SecretBoxAnimationAssets struct {
	Type             int    `json:"type"`
	BackgroundAsset  string `json:"background_asset"`
	AdditionalAsset1 string `json:"additional_asset_1"`
	AdditionalAsset2 string `json:"additional_asset_2"`
	AdditionalAsset3 string `json:"additional_asset_3"`
}

type SecretBoxCost struct {
	ID        int  `json:"id"`
	Type      int  `json:"type"`
	Amount    int  `json:"amount"`
	ItemID    *int `json:"item_id"`
	UnitCount int  `json:"unit_count"`
	Payable   bool `json:"payable"`
}

type SecretBoxButton struct {
	SecretBoxButtonType     int             `json:"secret_box_button_type"`
	CostList                []SecretBoxCost `json:"cost_list"`
	SecretBoxName           string          `json:"secret_box_name"`
	LimitCount              int             `json:"limit_count,omitempty"`
	LimitUnderMessage       string          `json:"limit_under_message,omitempty"`
	LimitUnderMessageSuffix string          `json:"limit_under_message_suffix,omitempty"`
	BalloonAsset            string          `json:"balloon_asset,omitempty"`
}

type SecretBoxInfo struct {
	SecretBoxID       int    `json:"secret_box_id"`
	SecretBoxType     int    `json:"secret_box_type"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	ShowEndDate       string `json:"show_end_date"`
	AddGauge          int    `json:"add_gauge"`
	AlwaysDisplayFlag int    `json:"always_display_flag"`
	PonCount          int    `json:"pon_count"`
	PonUpperLimit     int    `json:"pon_upper_limit"`
}

type SecretBoxPage struct {
	MenuAsset       string                   `json:"menu_asset"`
	PageOrder       int                      `json:"page_order"`
	AnimationAssets SecretBoxAnimationAssets `json:"animation_assets"`
	ButtonList      []SecretBoxButton        `json:"button_list"`
	SecretBoxInfo   SecretBoxInfo            `json:"secret_box_info"`
}

type MemberCategoryPageList struct {
	MemberCategory int             `json:"member_category"`
	PageList       []SecretBoxPage `json:"page_list"`
}

type AllData struct {
	UseCache           int                      `json:"use_cache"`
	IsUnitMax          bool                     `json:"is_unit_max"`
	ItemList           []ItemAmount             `json:"item_list"`
	GaugeInfo          GaugeInfo                `json:"gauge_info"`
	MemberCategoryList []MemberCategoryPageList `json:"member_category_list"`
}

type AllResp struct {
	Result     AllData `json:"result"`
	Status     int     `json:"status"`
	CommandNum bool    `json:"commandNum"`
	TimeStamp  int64   `json:"timeStamp"`
}
