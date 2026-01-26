package live

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func gameOver(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer func() {
		ss.ClearLiveInProgress()
		ss.Finalize()
	}()

	ss.Respond(liveschema.GameOverResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/live/gameover", middleware.Common, gameOver)
}
