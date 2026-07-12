package unit

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	unitschema "honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func favorite(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	var req unitschema.FavoriteReq
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	data := usermodel.UserUnitData{
		FavoriteFlag: req.FavoriteFlag != 0,
	}
	_, err = ss.UserEng.Table(new(usermodel.UserUnitData)).
		Where("user_id = ? AND unit_owning_user_id = ?", ss.UserID, req.UnitOwningUserID).
		Cols("favorite_flag").Update(&data)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(unitschema.FavoriteResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/favorite", middleware.Common, favorite)
}
