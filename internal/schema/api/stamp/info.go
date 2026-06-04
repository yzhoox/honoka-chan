package stampapischema

type StampList struct {
	Position int `json:"position"`
	StampID  int `json:"stamp_id"`
}

type SettingList struct {
	StampSettingID int         `json:"stamp_setting_id"`
	MainFlag       int         `json:"main_flag"`
	StampList      []StampList `json:"stamp_list"`
}

type StampSetting struct {
	StampType   int           `json:"stamp_type"`
	SettingList []SettingList `json:"setting_list"`
}

type InfoData struct {
	OwningStampIds []int          `json:"owning_stamp_ids"`
	StampSetting   []StampSetting `json:"stamp_setting"`
}

type InfoResp struct {
	Result     InfoData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
