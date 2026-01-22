package live

import (
	"honoka-chan/internal/schema/api/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func liveSchedule(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var liveDifficultyID []int
	specialLives := []live.SpecialLiveStatusList{}
	err = ss.MainEng.Table("special_live_m").Cols("live_difficulty_id").OrderBy("live_difficulty_id ASC").Find(&liveDifficultyID)
	if ss.CheckErr(err) {
		return
	}
	for _, id := range liveDifficultyID {
		specialLive := live.SpecialLiveStatusList{
			LiveDifficultyID:   id,
			Status:             1,
			HiScore:            0,
			HiComboCount:       0,
			ClearCnt:           0,
			AchievedGoalIDList: []int{},
		}
		specialLives = append(specialLives, specialLive)
	}

	livesList := []live.LiveList{}
	for _, v := range specialLives {
		livesList = append(livesList, live.LiveList{
			LiveDifficultyID: v.LiveDifficultyID,
			StartDate:        "2023-01-01 00:00:00",
			EndDate:          "2037-01-01 00:00:00",
			IsRandom:         false,
		})
	}
	res = live.ScheduleResp{
		Result: live.ScheduleData{
			EventList:              []any{},
			LiveList:               livesList,
			LimitedBonusList:       []any{},
			LimitedBonusCommonList: []live.LimitedBonusCommonList{}, // 特效道具
			RandomLiveList:         []live.RandomLiveList{},         // 随机歌曲
			FreeLiveList:           []any{},
			TrainingLiveList:       []live.TrainingLiveList{}, // 挑战歌曲
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
