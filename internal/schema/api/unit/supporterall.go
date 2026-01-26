package unitapischema

type SupporterList struct {
	UnitID int `json:"unit_id"`
	Amount int `json:"amount"`
}

type SupporterAllData struct {
	UnitSupportList []SupporterList `json:"unit_support_list"`
}

type SupporterAllResp struct {
	Result     SupporterAllData `json:"result"`
	Status     int              `json:"status"`
	CommandNum bool             `json:"commandNum"`
	TimeStamp  int64            `json:"timeStamp"`
}
