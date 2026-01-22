package unit

import (
	"honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func unitRemovableSkillInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var skillEquipCount []unit.SkillEquipCount
	err = ss.UserEng.Table("user_unit_skill_equip").Where("user_id = ?", ss.UserID).Select("unit_removable_skill_id,COUNT(*) AS ct").
		GroupBy("unit_removable_skill_id").Find(&skillEquipCount)
	if ss.CheckErr(err) {
		return
	}

	var rmSkillIds []int
	err = ss.MainEng.Table("unit_removable_skill_m").Where("effect_range = 1").Cols("unit_removable_skill_id").Find(&rmSkillIds)
	if ss.CheckErr(err) {
		return
	}

	owingInfo := []unit.OwningInfo{}
	for _, id := range rmSkillIds {
		info := unit.OwningInfo{
			UnitRemovableSkillID: id,
			TotalAmount:          9,
			EquippedAmount:       0,
			InsertDate:           "2023-01-01 12:00:00",
		}
		for _, sk := range skillEquipCount {
			if id == sk.UnitRemovableSkillId {
				info.EquippedAmount = sk.Count
				break
			}
		}
		owingInfo = append(owingInfo, info)
	}

	var unitOwningIds []int
	err = ss.UserEng.Table("user_unit_skill_equip").Where("user_id = ?", ss.UserID).Cols("unit_owning_user_id").GroupBy("unit_owning_user_id").Find(&unitOwningIds)
	if ss.CheckErr(err) {
		return
	}

	equipInfo := map[int]any{}
	for _, v := range unitOwningIds {
		detail := []unit.SkillEquipDetail{}
		err = ss.UserEng.Table("user_unit_skill_equip").Where("user_id = ? AND unit_owning_user_id = ?", ss.UserID, v).
			Cols("unit_removable_skill_id").Find(&detail)
		if ss.CheckErr(err) {
			return
		}

		equipInfo[v] = unit.SkillEquipList{
			UnitOwningUserID: v,
			Detail:           detail,
		}
	}

	res = unit.RemovableSkillInfoResp{
		Result: unit.RemovableSkillInfoData{
			OwningInfo:    owingInfo,
			EquipmentInfo: equipInfo,
		}, // 宝石
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
