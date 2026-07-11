package gdpr

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	gdprschema "honoka-chan/internal/schema/gdpr"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func get(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(gdprschema.GetResp{
		ResponseData: gdprschema.GetData{
			EnableGdpr:      true,
			IsEea:           false,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/gdpr/get", middleware.Common, get)
}
