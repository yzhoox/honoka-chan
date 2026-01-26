package unit

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	unitschema "honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func favorite(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	var req unitschema.FavoriteReq
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &req)
	if ss.CheckErr(err) {
		return
	}

	data := usermodel.UserUnitData{
		FavoriteFlag: req.FavoriteFlag != 0,
	}
	_, err = ss.UserEng.Table(new(usermodel.UserUnitData)).
		ID(req.UnitOwningUserID).Cols("favorite_flag").Update(&data)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(unitschema.FavoriteResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/favorite", middleware.Common, favorite)
}
