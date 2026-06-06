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

func reward(ctx *gin.Context) {
	ss := session.Get(ctx)
	token := ""
	defer func() {
		if ss.UserEng != nil {
			deleteRandomLiveByToken(ss, token)
			ss.ClearLiveInProgress()
		}
		ss.Finalize()
	}()

	req := rliveschema.RewardReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}
	token = req.Token

	randomLive, err := getRandomLiveByToken(ss, req.Token)
	if ss.CheckErr(err) {
		return
	}

	hasProgress, _ := ss.GetLiveInProgress()
	if !hasProgress {
		ss.CheckErr(errors.New("live progress not found"))
		return
	}

	if req.LiveDifficultyID != 0 && req.LiveDifficultyID != randomLive.LiveDifficultyID {
		ss.CheckErr(errors.New("random live difficulty mismatch"))
		return
	}

	rewardResp, err := livehandler.BuildRewardResp(ss, req.ToLiveRewardReq(randomLive.LiveDifficultyID), true)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(rewardResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/rlive/reward", middleware.Common, reward)
}
