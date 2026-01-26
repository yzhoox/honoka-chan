package unit

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	unitschema "honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func wearAccessory(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := unitschema.WearAccessoryReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &req)
	if ss.CheckErr(err) {
		return
	}

	// 取下饰品
	for _, v := range req.Remove {
		// fmt.Println("Remove:", v.AccessoryOwningUserID, v.UnitOwningUserID)
		_, err := ss.UserEng.Table("user_accessory_wear").
			Where("accessory_owning_user_id = ? AND unit_owning_user_id = ? AND user_id = ?", v.AccessoryOwningUserID, v.UnitOwningUserID, ss.UserID).
			Delete()
		if ss.CheckErr(err) {
			return
		}
	}

	// 佩戴饰品
	for _, v := range req.Wear {
		// fmt.Println("Wear:", v.AccessoryOwningUserID, v.UnitOwningUserID)
		data := unitschema.WearAccessoryData{
			AccessoryOwningUserID: v.AccessoryOwningUserID,
			UnitOwningUserID:      v.UnitOwningUserID,
			UserID:                ss.UserID,
		}
		_, err := ss.UserEng.Table("user_accessory_wear").Insert(&data)
		if ss.CheckErr(err) {
			return
		}
	}

	ss.Respond(unitschema.WearAccessoryResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/wearAccessory", middleware.Common, wearAccessory)
}
