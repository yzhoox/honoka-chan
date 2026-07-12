package rlive

import (
	"errors"
	livehandler "honoka-chan/internal/handler/live"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	rliveschema "honoka-chan/internal/schema/rlive"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func continuee(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	req := rliveschema.TokenReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	_, err = getRandomLiveByToken(ss, req.Token)
	if ss.CheckErr(err) {
		return
	}

	hasProgress, _ := ss.GetLiveInProgress()
	if !hasProgress {
		ss.CheckErr(errors.New("live progress not found"))
		return
	}

	ss.Respond(livehandler.BuildContinueResp(ss))
}

func init() {
	router.AddHandler("main.php", "POST", "/rlive/continue", middleware.Common, continuee)
}
