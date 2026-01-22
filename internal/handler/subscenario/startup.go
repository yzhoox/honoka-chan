package subscenario

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/subscenario"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func startup(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	startReq := subscenario.StartUpReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &startReq)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(subscenario.StartUpResp{
		ResponseData: subscenario.StartUpData{
			SubscenarioID:      startReq.SubscenarioID,
			ScenarioAdjustment: 50,
			ServerTimestamp:    time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/subscenario/startup", middleware.Common, startup)
}
