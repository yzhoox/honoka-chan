package unitapischema

type AccessoryMaterialList struct {
	AccessoryID int `json:"accessory_id"`
	Amount      int `json:"amount"`
}

type AccessoryMaterialAllData struct {
	AccessoryMaterialList []AccessoryMaterialList `json:"accessory_material_list"`
}

type AccessoryMaterialAllResp struct {
	Result     AccessoryMaterialAllData `json:"result"`
	Status     int                      `json:"status"`
	CommandNum bool                     `json:"commandNum"`
	TimeStamp  int64                    `json:"timeStamp"`
}
