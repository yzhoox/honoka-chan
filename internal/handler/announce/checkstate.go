package announce

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	announceschema "honoka-chan/internal/schema/announce"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func checkState(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(announceschema.CheckStateDataResp{
		ResponseData: announceschema.CheckStateData{
			HasUnreadAnnounce: false,
			ServerTimestamp:   time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/announce/checkState", middleware.Common, checkState)
}
