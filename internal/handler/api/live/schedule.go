package live

import (
	liveapischema "honoka-chan/internal/schema/api/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func liveSchedule(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var liveID []int
	err = ss.MainEng.Table("special_live_m").
		Cols("live_difficulty_id").OrderBy("live_difficulty_id ASC").Find(&liveID)
	if ss.CheckErr(err) {
		return
	}

	liveList := []liveapischema.LiveList{}
	for _, id := range liveID {
		liveList = append(liveList, liveapischema.LiveList{
			LiveDifficultyID: id,
			StartDate:        "2023-01-01 00:00:00",
			EndDate:          "2037-01-01 00:00:00",
			IsRandom:         false,
		})
	}

	res = liveapischema.ScheduleResp{
		Result: liveapischema.ScheduleData{
			EventList:              []any{},
			LiveList:               liveList,
			LimitedBonusList:       []any{},
			LimitedBonusCommonList: []liveapischema.LimitedBonusCommonList{},
			RandomLiveList:         []liveapischema.RandomLiveList{},
			FreeLiveList:           []any{},
			TrainingLiveList:       []liveapischema.TrainingLiveList{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
