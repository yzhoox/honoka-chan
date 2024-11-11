package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TosCheck(ctx *gin.Context) {
	tosResp := model.TosResp{
		ResponseData: model.TosRes{
			TosID:           1,
			TosType:         1,
			IsAgreed:        true,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(tosResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
