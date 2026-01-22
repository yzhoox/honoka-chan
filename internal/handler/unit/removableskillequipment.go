package unit

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func removableSkillEquipment(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := unit.RemovableSkillEquipmentReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &req)
	if ss.CheckErr(err) {
		return
	}

	// 取下宝石
	for _, v := range req.Remove {
		// fmt.Println("Remove:", v.UnitOwningUserID, v.UnitRemovableSkillID)
		_, err := ss.UserEng.Table("user_unit_skill_equip").
			Where("unit_removable_skill_id = ? AND unit_owning_user_id = ? AND user_id = ?", v.UnitRemovableSkillID, v.UnitOwningUserID, ss.UserID).
			Delete()
		if ss.CheckErr(err) {
			return
		}
	}

	// 佩戴宝石
	for _, v := range req.Equip {
		// fmt.Println("Equip:", v.UnitOwningUserID, v.UnitRemovableSkillID)
		data := unit.RemovableSkillEquipmentData{
			UnitRemovableSkillId: v.UnitRemovableSkillID,
			UnitOwningUserID:     v.UnitOwningUserID,
			UserID:               ss.UserID,
		}
		_, err := ss.UserEng.Table("user_unit_skill_equip").Insert(&data)
		if ss.CheckErr(err) {
			return
		}
	}

	ss.Respond(unit.RemovableSkillEquipmentResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/removableSkillEquipment", middleware.Common, removableSkillEquipment)
}
