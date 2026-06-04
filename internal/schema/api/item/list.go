package itemapischema

type GeneralItemList struct {
	ItemID          int  `json:"item_id"`
	Amount          int  `json:"amount"`
	UseButtonFlag   bool `json:"use_button_flag"`
	GeneralItemType int  `json:"general_item_type"`
}

type BuffItemList struct {
	ItemID   int `json:"item_id"`
	Amount   int `json:"amount"`
	BuffType int `json:"buff_type"`
}

type ListData struct {
	GeneralItemList []GeneralItemList `json:"general_item_list"`
	BuffItemList    []BuffItemList    `json:"buff_item_list"`
}

type ListResp struct {
	Result     ListData `json:"result"`
	Status     int      `json:"status"`
	CommandNum bool     `json:"commandNum"`
	TimeStamp  int64    `json:"timeStamp"`
}
