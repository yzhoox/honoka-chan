package live

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func continuee(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(BuildContinueResp(ss))
}

func BuildContinueResp(ss *session.Session) liveschema.ContinueResp {
	return liveschema.ContinueResp{
		ResponseData: liveschema.ContinueData{
			BeforeSnsCoin:   ss.UserPref.SnsCoin,
			AfterSnsCoin:    ss.UserPref.SnsCoin,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	}
}

func init() {
	router.AddHandler("main.php", "POST", "/live/continue", middleware.Common, continuee)
}
