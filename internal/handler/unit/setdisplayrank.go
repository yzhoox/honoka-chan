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

func setDisplayRank(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	var req unitschema.SetDisplayRankReq
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	data := usermodel.UserUnitData{
		DisplayRank: req.DisplayRank,
	}
	_, err = ss.UserEng.Table(new(usermodel.UserUnitData)).
		ID(req.UnitOwningUserID).Cols("display_rank").Update(&data)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(unitschema.SetDisplayRankResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/setDisplayRank", middleware.Common, setDisplayRank)
}
