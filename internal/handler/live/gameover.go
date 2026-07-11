package live

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func gameOver(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer func() {
		if ss.UserEng != nil {
			ss.ClearLiveInProgress()
		}
		ss.Finalize()
	}()

	ss.Respond(BuildGameOverResp())
}

func BuildGameOverResp() liveschema.GameOverResp {
	return liveschema.GameOverResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	}
}

func init() {
	router.AddHandler("main.php", "POST", "/live/gameover", middleware.Common, gameOver)
}
