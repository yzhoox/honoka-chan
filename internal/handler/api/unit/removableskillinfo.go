package unit

import (
	unitapischema "honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func unitRemovableSkillInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var skillEquipCount []unitapischema.SkillEquipCount
	err = ss.UserEng.Table("user_unit_skill_equip").Where("user_id = ?", ss.UserID).Select("unit_removable_skill_id,COUNT(*) AS ct").
		GroupBy("unit_removable_skill_id").Find(&skillEquipCount)
	if err != nil {
		return nil, err
	}

	var rmSkillIds []int
	err = ss.MainEng.Table("unit_removable_skill_m").Where("effect_range = 1").Cols("unit_removable_skill_id").Find(&rmSkillIds)
	if err != nil {
		return nil, err
	}

	owingInfo := []unitapischema.OwningInfo{}
	for _, id := range rmSkillIds {
		info := unitapischema.OwningInfo{
			UnitRemovableSkillID: id,
			TotalAmount:          999,
			EquippedAmount:       0, // TODO: 计算实际数量
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
	if err != nil {
		return nil, err
	}

	equipInfo := map[int]any{}
	for _, v := range unitOwningIds {
		detail := []unitapischema.SkillEquipDetail{}
		err = ss.UserEng.Table("user_unit_skill_equip").Where("user_id = ? AND unit_owning_user_id = ?", ss.UserID, v).
			Cols("unit_removable_skill_id").Find(&detail)
		if err != nil {
			return nil, err
		}

		equipInfo[v] = unitapischema.SkillEquipList{
			UnitOwningUserID: v,
			Detail:           detail,
		}
	}

	res = unitapischema.RemovableSkillInfoResp{
		Result: unitapischema.RemovableSkillInfoData{
			OwningInfo:    owingInfo,
			EquipmentInfo: equipInfo,
		}, // 宝石
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
