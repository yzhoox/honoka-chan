package unit

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/unit"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func setDisplayRank(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(unit.SetDisplayRankResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/setDisplayRank", middleware.Common, setDisplayRank)
}
