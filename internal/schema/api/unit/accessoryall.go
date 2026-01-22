package unit

type AccessoryList struct {
	AccessoryOwningUserID int  `json:"accessory_owning_user_id" xorm:"accessory_owning_user_id"`
	AccessoryID           int  `json:"accessory_id" xorm:"accessory_id"`
	Exp                   int  `json:"exp" xorm:"exp"`
	NextExp               int  `json:"next_exp" xorm:"-"`
	Level                 int  `json:"level" xorm:"-"`
	MaxLevel              int  `json:"max_level" xorm:"-"`
	RankUpCount           int  `json:"rank_up_count" xorm:"-"`
	FavoriteFlag          bool `json:"favorite_flag" xorm:"-"`
}

type WearingInfo struct {
	UnitOwningUserID      int `json:"unit_owning_user_id" xorm:"unit_owning_user_id"`
	AccessoryOwningUserID int `json:"accessory_owning_user_id" xorm:"accessory_owning_user_id"`
}

type AccessoryAllData struct {
	AccessoryList      []AccessoryList `json:"accessory_list"`
	WearingInfo        []WearingInfo   `json:"wearing_info"`
	EspecialCreateFlag bool            `json:"especial_create_flag"`
}

type AccessoryAllResp struct {
	Result     AccessoryAllData `json:"result"`
	Status     int              `json:"status"`
	CommandNum bool             `json:"commandNum"`
	TimeStamp  int64            `json:"timeStamp"`
}
