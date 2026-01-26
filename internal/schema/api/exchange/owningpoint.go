package exchangeapischema

type ExchangePointList struct {
	Rarity        int `json:"rarity"`
	ExchangePoint int `json:"exchange_point"`
}

type OwningPointData struct {
	ExchangePointList []ExchangePointList `json:"exchange_point_list"`
}

type OwningPointResp struct {
	Result     OwningPointData `json:"result"`
	Status     int             `json:"status"`
	CommandNum bool            `json:"commandNum"`
	TimeStamp  int64           `json:"timeStamp"`
}
