package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Gdpr(ctx *gin.Context) {
	gdprResp := model.GdprResp{
		ResponseData: model.GdprRes{
			EnableGdpr:      true,
			IsEea:           false,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(gdprResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
