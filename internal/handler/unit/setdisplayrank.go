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

func setDisplayRank(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	var req unitschema.SetDisplayRankReq
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &req)
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
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/setDisplayRank", middleware.Common, setDisplayRank)
}
