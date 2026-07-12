package reward

import (
	"honoka-chan/internal/constant"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	commonschema "honoka-chan/internal/schema/common"
	rewardschema "honoka-chan/internal/schema/reward"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func sellUnit(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	ss.Respond(rewardschema.SellUnitResp{
		ResponseData: commonschema.ErrorData{
			ErrorCode: constant.ErrorCodeNoUnitIsSellable,
		},
		ReleaseInfo: []any{},
		StatusCode:  constant.ErrorCodeAcceptableError,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/reward/sellUnit", middleware.Common, sellUnit)
}
