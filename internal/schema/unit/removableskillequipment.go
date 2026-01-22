package unit

type SkillRemove struct {
	UnitRemovableSkillID int `json:"unit_removable_skill_id"`
	UnitOwningUserID     int `json:"unit_owning_user_id"`
}

type SkillEquip struct {
	UnitRemovableSkillID int `json:"unit_removable_skill_id"`
	UnitOwningUserID     int `json:"unit_owning_user_id"`
}

type RemovableSkillEquipmentReq struct {
	Module     string        `json:"module"`
	Remove     []SkillRemove `json:"remove"`
	Action     string        `json:"action"`
	TimeStamp  int           `json:"timeStamp"`
	Equip      []SkillEquip  `json:"equip"`
	Mgd        int           `json:"mgd"`
	CommandNum string        `json:"commandNum"`
}

type RemovableSkillEquipmentData struct {
	Id                   int `xorm:"id pk autoincr"`
	UnitRemovableSkillId int `xorm:"unit_removable_skill_id"`
	UnitOwningUserID     int `xorm:"unit_owning_user_id"`
	UserID               int `xorm:"user_id"`
}

type RemovableSkillEquipmentResp struct {
	ResponseData []any `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}
