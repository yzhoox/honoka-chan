package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func TosCheck(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

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

	ss.Respond(tosResp)
}
