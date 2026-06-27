package session

import (
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"

	"xorm.io/xorm"
)

func (ss *Session) GetBasicUnitInfo() *xorm.Session {
	return ss.UserEng.Table("user_unit_data").Alias("a").
		Join("LEFT", "common_unit_data", "a.unit_id = common_unit_data.unit_id").
		Cols(`
            a.unit_owning_user_id,
            a.favorite_flag,
            a.display_rank,
            common_unit_data.*
        `)
}

func (ss *Session) GetUnitInfo(unitID int) (bool, *unitmodel.CommonUnitData) {
	if ss.UserEng == nil {
		return false, nil
	}

	unitInfo := unitmodel.CommonUnitData{}
	has, err := ss.UserEng.Table(new(unitmodel.CommonUnitData)).
		Where("unit_id = ?", unitID).Get(&unitInfo)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &unitInfo
}

func (ss *Session) GetUserUnitInfo(unitOwningUserID int) (bool, *unitmodel.CommonUnitData) {
	var unitID int
	has, err := ss.UserEng.Table(new(usermodel.UserUnitData)).
		Where("unit_owning_user_id = ?", unitOwningUserID).
		Cols("unit_id").Get(&unitID)
	if ss.CheckErr(err) {
		return false, nil
	}

	if !has {
		return false, nil
	}

	return ss.GetUnitInfo(unitID)
}

func (ss *Session) GetUserUnitSkillEquip(unitOwningUserID int) []usermodel.UserUnitSkillEquip {
	skill := []usermodel.UserUnitSkillEquip{}
	err := ss.UserEng.Table("user_unit_skill_equip").
		Where("unit_owning_user_id = ?", unitOwningUserID).Find(&skill)
	if ss.CheckErr(err) {
		return []usermodel.UserUnitSkillEquip{}
	}

	return skill
}

func (ss *Session) GetUserUnitSkillEquipID(unitOwningUserID int) []int {
	skillID := []int{}
	skill := ss.GetUserUnitSkillEquip(unitOwningUserID)
	for _, s := range skill {
		skillID = append(skillID, s.UnitRemovableSkillID)
	}

	return skillID
}
