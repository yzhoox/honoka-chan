package gdpr

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/gdpr"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func get(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(gdpr.GetResp{
		ResponseData: gdpr.GetData{
			EnableGdpr:      true,
			IsEea:           false,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/gdpr/get", middleware.Common, get)
}
