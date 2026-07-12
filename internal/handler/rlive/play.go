package rlive

import (
	livehandler "honoka-chan/internal/handler/live"
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	rliveschema "honoka-chan/internal/schema/rlive"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func play(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	req := rliveschema.PlayReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	randomLive, err := getRandomLiveByToken(ss, req.Token)
	if ss.CheckErr(err) {
		return
	}

	if err := ss.ResetRandomLiveInProgress(); ss.CheckErr(err) {
		return
	}

	ss.ClearLiveInProgress()
	ss.RegisterLiveInProgress(req.UnitDeckID)

	_, err = ss.UserEng.Table(new(usermodel.UserLiveRandom)).
		Where("user_id = ? AND token = ?", ss.UserID, req.Token).
		Cols("in_progress").
		Update(&usermodel.UserLiveRandom{InProgress: true})
	if ss.CheckErr(err) {
		return
	}

	playResp, err := livehandler.BuildPlayResp(ss, req.ToLivePlayReq(randomLive.LiveDifficultyID), true)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(playResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/rlive/play", middleware.Common, play)
}
