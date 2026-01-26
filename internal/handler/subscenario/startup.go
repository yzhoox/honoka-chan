package subscenario

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	subscenarioschema "honoka-chan/internal/schema/subscenario"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func startup(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	startReq := subscenarioschema.StartUpReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &startReq)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(subscenarioschema.StartUpResp{
		ResponseData: subscenarioschema.StartUpData{
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
