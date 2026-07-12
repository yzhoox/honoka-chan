package unit

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	unitschema "honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func sale(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	req := unitschema.SaleReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	totalCoin := 0
	saleDetail := []unitschema.Detail{}
	for _, id := range req.UnitOwningUserID {
		salePrice := 1 // TODO: 从 unit_level_up_pattern_m 表获取实际价格
		totalCoin += salePrice

		_, unitData := ss.GetUserUnitInfo(id)

		saleDetail = append(saleDetail, unitschema.Detail{
			UnitOwningUserID: id,
			UnitID:           unitData.UnitID,
			IsSigned:         unitData.IsSigned,
			Price:            salePrice,
		})

		// 卸载宝石
		_, err = ss.UserEng.Table(new(usermodel.UserUnitSkillEquip)).
			Where("unit_owning_user_id = ?", id).Delete()
		if ss.CheckErr(err) {
			return
		}

		// 卸载饰品
		_, err = ss.UserEng.Table(new(usermodel.UserAccessoryWear)).
			Where("unit_owning_user_id = ?", id).Delete()
		if ss.CheckErr(err) {
			return
		}

		// 移除卡片
		_, err = ss.UserEng.Table(new(usermodel.UserUnitData)).
			Where("unit_owning_user_id = ?", id).Delete()
		if ss.CheckErr(err) {
			return
		}
	}

	skillList := []int{}
	err = ss.MainEng.Table("unit_removable_skill_m").
		Where("skill_type = 1").Cols("unit_removable_skill_id").Find(&skillList)
	if ss.CheckErr(err) {
		return
	}

	owningData := []unitschema.OwningInfo{}
	for _, sk := range skillList {
		owningData = append(owningData, unitschema.OwningInfo{
			UnitRemovableSkillID: sk,
			TotalAmount:          999,
			EquippedAmount:       0, // TODO: 计算实际数量
			InsertDate:           "2023-01-01 12:00:00",
		})
	}

	ss.Respond(unitschema.SaleResp{
		ResponseData: unitschema.SaleData{
			Total:                totalCoin,
			Detail:               saleDetail,
			BeforeUserInfo:       ss.GetUserInfo(),
			AfterUserInfo:        ss.GetUserInfo(),
			RewardBoxFlag:        false,
			GetExchangePointList: []any{},
			UnitRemovableSkill: unitschema.UnitRemovableSkill{
				OwningInfo: owningData,
			},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/sale", middleware.Common, sale)
}
