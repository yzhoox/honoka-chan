package reward

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	rewardschema "honoka-chan/internal/schema/reward"
	"honoka-chan/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func list(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(rewardschema.ListResp{
		ResponseData: rewardschema.ListData{
			ItemCount: 0,
			Limit:     20,
			Order:     0,
			Items:     []any{},
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/reward/rewardList", middleware.Common, list)
}
