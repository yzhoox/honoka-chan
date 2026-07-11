package multiunit

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	multiunitschema "honoka-chan/internal/schema/multiunit"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func scenarioStartup(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	startReq := multiunitschema.ScenarioStartUpReq{}
	err := honokautils.ParseRequestData(ctx, &startReq)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(multiunitschema.ScenarioStartUpResp{
		ResponseData: multiunitschema.ScenarioStartUpData{
			MultiUnitScenarioID: startReq.MultiUnitScenarioID,
			ScenarioAdjustment:  50,
			ServerTimestamp:     time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/multiunit/scenarioStartup", middleware.Common, scenarioStartup)
}
