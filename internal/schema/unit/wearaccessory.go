package unit

type AccessoryRemove struct {
	AccessoryOwningUserID int `json:"accessory_owning_user_id"`
	UnitOwningUserID      int `json:"unit_owning_user_id"`
}

type AccessoryWear struct {
	AccessoryOwningUserID int `json:"accessory_owning_user_id"`
	UnitOwningUserID      int `json:"unit_owning_user_id"`
}

type WearAccessoryReq struct {
	Module     string            `json:"module"`
	Remove     []AccessoryRemove `json:"remove"`
	Action     string            `json:"action"`
	TimeStamp  int               `json:"timeStamp"`
	Wear       []AccessoryWear   `json:"wear"`
	Mgd        int               `json:"mgd"`
	CommandNum string            `json:"commandNum"`
}

type WearAccessoryData struct {
	Id                    int `xorm:"id pk autoincr"`
	AccessoryOwningUserID int `xorm:"accessory_owning_user_id"`
	UnitOwningUserID      int `xorm:"unit_owning_user_id"`
	UserID                int `xorm:"user_id"`
}

type WearAccessoryResp struct {
	ResponseData []any `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}
