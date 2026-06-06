package live

import (
	liveapischema "honoka-chan/internal/schema/api/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func liveSchedule(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	now := time.Now()
	liveIDs, err := listTodaySpecialRotationDifficultyIDs(ss, now)
	if err != nil {
		return nil, err
	}

	liveList := make([]liveapischema.LiveList, 0, len(liveIDs))
	startDate := startOfDay(now.In(jst)).Format("2006-01-02 15:04:05")
	endDate := nextDayStart(now.In(jst)).Format("2006-01-02 15:04:05")
	for _, id := range liveIDs {
		liveList = append(liveList, liveapischema.LiveList{
			LiveDifficultyID: id,
			StartDate:        startDate,
			EndDate:          endDate,
			IsRandom:         false,
		})
	}

	res = liveapischema.ScheduleResp{
		Result: liveapischema.ScheduleData{
			EventList:              []any{},
			LiveList:               liveList,
			LimitedBonusList:       []any{},
			LimitedBonusCommonList: []liveapischema.LimitedBonusCommonList{},
			RandomLiveList:         buildRandomLiveList(),
			FreeLiveList:           []any{},
			TrainingLiveList:       []liveapischema.TrainingLiveList{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
