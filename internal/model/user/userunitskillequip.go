package usermodel

type UserUnitSkillEquip struct {
	ID                   int `xorm:"id pk autoincr"`
	UnitOwningUserID     int `xorm:"unit_owning_user_id"`
	UnitRemovableSkillID int `xorm:"unit_removable_skill_id"`
	UserID               int `xorm:"user_id"`
}

func (UserUnitSkillEquip) TableName() string {
	return "user_unit_skill_equip"
}
