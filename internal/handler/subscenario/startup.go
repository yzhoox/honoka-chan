package subscenario

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	subscenarioschema "honoka-chan/internal/schema/subscenario"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func startup(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	startReq := subscenarioschema.StartUpReq{}
	err := honokautils.ParseRequestData(ctx, &startReq)
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
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/subscenario/startup", middleware.Common, startup)
}
