package handler

import (
	"encoding/json"
	"honoka-chan/internal/model"
	"honoka-chan/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func MultiUnitStartUp(ctx *gin.Context) {
	startReq := model.MultiUnitStartUpReq{}
	err := json.Unmarshal([]byte(ctx.PostForm("request_data")), &startReq)
	utils.CheckErr(err)

	startResp := model.MultiUnitStartUpResp{
		ResponseData: model.MultiUnitStartUpRes{
			MultiUnitScenarioID: startReq.MultiUnitScenarioID,
			ScenarioAdjustment:  50,
			ServerTimestamp:     time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(startResp)
	utils.CheckErr(err)

	ctx.Header("X-Message-Sign", utils.GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
