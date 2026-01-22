package announce

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/announce"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func checkState(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(announce.CheckStateDataResp{
		ResponseData: announce.CheckStateData{
			HasUnreadAnnounce: false,
			ServerTimestamp:   time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/announce/checkState", middleware.Common, checkState)
}
