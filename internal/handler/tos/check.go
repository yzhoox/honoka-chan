package tos

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	tosschema "honoka-chan/internal/schema/tos"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func check(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(tosschema.CheckResp{
		ResponseData: tosschema.CheckData{
			TosID:           1,
			TosType:         1,
			IsAgreed:        true,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/tos/tosCheck", middleware.Common, check)
}
