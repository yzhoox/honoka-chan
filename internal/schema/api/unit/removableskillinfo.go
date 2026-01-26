package unitapischema

type SkillEquipCount struct {
	UnitRemovableSkillId int `xorm:"unit_removable_skill_id"`
	Count                int `xorm:"ct"`
}

type SkillEquipData struct {
	Id                   int    `xorm:"id pk autoincr"`
	UnitRemovableSkillId int    `xorm:"unit_removable_skill_id"`
	UnitOwningUserID     int    `xorm:"unit_owning_user_id"`
	UserId               string `xorm:"user_id"`
}

type SkillEquipDetail struct {
	UnitRemovableSkillID int `json:"unit_removable_skill_id" xorm:"unit_removable_skill_id"`
}

type SkillEquipList struct {
	UnitOwningUserID int                `json:"unit_owning_user_id"`
	Detail           []SkillEquipDetail `json:"detail"`
}

type OwningInfo struct {
	UnitRemovableSkillID int    `json:"unit_removable_skill_id"`
	TotalAmount          int    `json:"total_amount"`
	EquippedAmount       int    `json:"equipped_amount"`
	InsertDate           string `json:"insert_date"`
}

type RemovableSkillInfoData struct {
	OwningInfo    []OwningInfo `json:"owning_info"`
	EquipmentInfo map[int]any  `json:"equipment_info"`
}

type RemovableSkillInfoResp struct {
	Result     RemovableSkillInfoData `json:"result"`
	Status     int                    `json:"status"`
	CommandNum bool                   `json:"commandNum"`
	TimeStamp  int64                  `json:"timeStamp"`
}
