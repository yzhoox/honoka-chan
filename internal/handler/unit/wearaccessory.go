package unit

import (
	"errors"
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	unitschema "honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func wearAccessory(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	req := unitschema.WearAccessoryReq{}
	err := honokautils.ParseRequestData(ctx, &req)
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
	var exists bool
	for _, v := range req.Wear {
		// fmt.Println("Wear:", v.AccessoryOwningUserID, v.UnitOwningUserID)
		exists, err = ss.UserEng.Table(new(usermodel.UserAccessory)).
			Where("user_id = ? AND accessory_owning_user_id = ?", ss.UserID, v.AccessoryOwningUserID).
			Exist()
		if ss.CheckErr(err) {
			return
		}
		if !exists {
			ss.CheckErr(errors.New("accessory not found for user"))
			return
		}

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
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/wearAccessory", middleware.Common, wearAccessory)
}
