package live

import (
	liveapischema "honoka-chan/internal/schema/api/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func liveStatus(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var liveDifficultyID []int
	normalLives := []liveapischema.NormalLiveStatusList{}
	err = ss.MainEng.Table("normal_live_m").Cols("live_difficulty_id").OrderBy("live_difficulty_id ASC").Find(&liveDifficultyID)
	if ss.CheckErr(err) {
		return
	}
	for _, id := range liveDifficultyID {
		normalLive := liveapischema.NormalLiveStatusList{
			LiveDifficultyID:   id,
			Status:             1,
			HiScore:            0,
			HiComboCount:       0,
			ClearCnt:           0,
			AchievedGoalIDList: []int{},
		}
		normalLives = append(normalLives, normalLive)
	}

	specialLives := []liveapischema.SpecialLiveStatusList{}
	err = ss.MainEng.Table("special_live_m").Cols("live_difficulty_id").OrderBy("live_difficulty_id ASC").Find(&liveDifficultyID)
	if ss.CheckErr(err) {
		return
	}
	for _, id := range liveDifficultyID {
		specialLive := liveapischema.SpecialLiveStatusList{
			LiveDifficultyID:   id,
			Status:             1,
			HiScore:            0,
			HiComboCount:       0,
			ClearCnt:           0,
			AchievedGoalIDList: []int{},
		}
		specialLives = append(specialLives, specialLive)
	}

	res = liveapischema.StatusResp{
		Result: liveapischema.StatusData{
			NormalLiveStatusList:   normalLives,
			SpecialLiveStatusList:  specialLives,
			TrainingLiveStatusList: []liveapischema.TrainingLiveStatusList{},
			MarathonLiveStatusList: []any{},
			FreeLiveStatusList:     []any{},
			CanResumeLive:          false,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
