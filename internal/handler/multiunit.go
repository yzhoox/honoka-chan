package handler

import (
	"encoding/json"
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func MultiUnitStartUp(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	startReq := model.MultiUnitStartUpReq{}
	err := json.Unmarshal([]byte(ctx.PostForm("request_data")), &startReq)
	if ss.CheckErr(err) {
		return
	}

	startResp := model.MultiUnitStartUpResp{
		ResponseData: model.MultiUnitStartUpRes{
			MultiUnitScenarioID: startReq.MultiUnitScenarioID,
			ScenarioAdjustment:  50,
			ServerTimestamp:     time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(startResp)
}
