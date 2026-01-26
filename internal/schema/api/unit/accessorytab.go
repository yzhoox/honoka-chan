package unitapischema

type AssetList struct {
	UnitTypeID int    `json:"unit_type_id"`
	AssetPath  string `json:"asset_path"`
}

type TabList struct {
	TabName   string      `json:"tab_name"`
	AssetList []AssetList `json:"asset_list"`
}

type AccessoryTabData struct {
	TabList []TabList `json:"tab_list"`
}

type AccessoryTabResp struct {
	Result     AccessoryTabData `json:"result"`
	Status     int              `json:"status"`
	CommandNum bool             `json:"commandNum"`
	TimeStamp  int64            `json:"timeStamp"`
}
