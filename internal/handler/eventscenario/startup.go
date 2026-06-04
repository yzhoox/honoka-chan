package eventscenario

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	eventscenarioschema "honoka-chan/internal/schema/eventscenario"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func startup(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	startReq := eventscenarioschema.StartUpReq{}
	err := honokautils.ParseRequestData(ctx, &startReq)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(eventscenarioschema.StartUpResp{
		ResponseData: eventscenarioschema.StartUpData{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/eventscenario/startup", middleware.Common, startup)
}
