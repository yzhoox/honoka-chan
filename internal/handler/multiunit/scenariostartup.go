package multiunit

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	multiunitschema "honoka-chan/internal/schema/multiunit"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func scenarioStartup(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	startReq := multiunitschema.ScenarioStartUpReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &startReq)
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
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/multiunit/scenarioStartup", middleware.Common, scenarioStartup)
}
