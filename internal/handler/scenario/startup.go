package scenario

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/scenario"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func startup(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	startReq := scenario.StartUpReq{}
	err := json.Unmarshal([]byte(ctx.PostForm("request_data")), &startReq)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(scenario.StartUpResp{
		ResponseData: scenario.StartUpData{
			ScenarioID:         startReq.ScenarioID,
			ScenarioAdjustment: 50,
			ServerTimestamp:    time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/scenario/startup", middleware.Common, startup)
}
