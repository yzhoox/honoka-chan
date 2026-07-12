package rlive

import (
	livehandler "honoka-chan/internal/handler/live"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	rliveschema "honoka-chan/internal/schema/rlive"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func gameOver(ctx *gin.Context) {
	ss := session.Get(ctx)
	token := ""
	defer ss.FinalizeOrRollbackAfter(func() {
		if ss.UserEng != nil {
			deleteRandomLiveByToken(ss, token)
			ss.ClearLiveInProgress()
		}
	})

	req := rliveschema.TokenReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}
	token = req.Token

	_, err = getRandomLiveByToken(ss, req.Token)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(livehandler.BuildGameOverResp())
}

func init() {
	router.AddHandler("main.php", "POST", "/rlive/gameover", middleware.Common, gameOver)
}
