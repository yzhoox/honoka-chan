package handler

import (
	"encoding/json"
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func SubScenarioStartup(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	startReq := model.SubScenarioReq{}
	err := json.Unmarshal([]byte(ctx.PostForm("request_data")), &startReq)
	if ss.CheckErr(err) {
		return
	}

	startResp := model.SubScenarioResp{
		ResponseData: model.SubScenarioRes{
			SubscenarioID:      startReq.SubscenarioID,
			ScenarioAdjustment: 50,
			ServerTimestamp:    time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(startResp)
}

func SubScenarioReward(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	data := honokautils.ReadAllText("assets/serverdata/subreward.json")
	var resp map[string]any
	err := json.Unmarshal([]byte(data), &resp)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(resp)
}
