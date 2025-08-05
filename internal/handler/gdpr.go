package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func Gdpr(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	gdprResp := model.GdprResp{
		ResponseData: model.GdprRes{
			EnableGdpr:      true,
			IsEea:           false,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(gdprResp)
}
