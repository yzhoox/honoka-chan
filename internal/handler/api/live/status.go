package live

import (
	liveapischema "honoka-chan/internal/schema/api/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func liveStatus(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	normalIDs, err := listAvailableNormalLiveDifficultyIDs(ss)
	if err != nil {
		return nil, err
	}

	specialIDs, err := listTodaySpecialRotationDifficultyIDs(ss, time.Now())
	if err != nil {
		return nil, err
	}
	trainingIDs, err := listAvailableTrainingLiveDifficultyIDs(ss)
	if err != nil {
		return nil, err
	}

	res = liveapischema.StatusResp{
		Result: liveapischema.StatusData{
			NormalLiveStatusList:   buildNormalLiveStatusList(normalIDs),
			SpecialLiveStatusList:  buildSpecialLiveStatusList(specialIDs),
			TrainingLiveStatusList: buildTrainingLiveStatusList(trainingIDs),
			MarathonLiveStatusList: []any{},
			FreeLiveStatusList:     []any{},
			CanResumeLive:          true,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
