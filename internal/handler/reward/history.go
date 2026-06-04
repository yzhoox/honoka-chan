package reward

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	rewardschema "honoka-chan/internal/schema/reward"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func history(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(rewardschema.HistoryResp{
		ResponseData: rewardschema.HistoryData{
			ItemCount: 0,
			Limit:     20,
			History:   []any{},
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/reward/rewardHistory", middleware.Common, history)
}
