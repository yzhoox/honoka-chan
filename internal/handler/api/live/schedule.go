package live

import (
	liveapischema "honoka-chan/internal/schema/api/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func liveSchedule(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	schedules, err := listCurrentAndNextSpecialRotationSchedules(ss, time.Now())
	if err != nil {
		return nil, err
	}

	liveList := make([]liveapischema.LiveList, 0, len(schedules))
	for _, schedule := range schedules {
		liveList = append(liveList, liveapischema.LiveList{
			LiveDifficultyID: schedule.LiveDifficultyID,
			StartDate:        schedule.StartTime.Format("2006-01-02 15:04:05"),
			EndDate:          schedule.EndTime.Format("2006-01-02 15:04:05"),
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
