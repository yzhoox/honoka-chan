package live

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func continuee(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(liveschema.ContinueResp{
		ResponseData: liveschema.ContinueData{
			BeforeSnsCoin:   ss.UserPref.SnsCoin,
			AfterSnsCoin:    ss.UserPref.SnsCoin,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/live/continue", middleware.Common, continuee)
}
