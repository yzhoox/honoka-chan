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

	specialSchedules, err := listCurrentAndNextSpecialRotationSchedules(ss, time.Now())
	if err != nil {
		return nil, err
	}
	specialIDs := make([]int, 0, len(specialSchedules))
	for _, schedule := range specialSchedules {
		specialIDs = append(specialIDs, schedule.LiveDifficultyID)
	}
	trainingIDs, err := listAvailableTrainingLiveDifficultyIDs(ss)
	if err != nil {
		return nil, err
	}
	allIDs := make([]int, 0, len(normalIDs)+len(specialIDs)+len(trainingIDs))
	allIDs = append(allIDs, normalIDs...)
	allIDs = append(allIDs, specialIDs...)
	allIDs = append(allIDs, trainingIDs...)
	statusSnapshotMap, err := ss.BuildLiveStatusSnapshotMap(allIDs)
	if err != nil {
		return nil, err
	}

	res = liveapischema.StatusResp{
		Result: liveapischema.StatusData{
			NormalLiveStatusList:   buildNormalLiveStatusList(normalIDs, statusSnapshotMap),
			SpecialLiveStatusList:  buildSpecialLiveStatusList(specialIDs, statusSnapshotMap),
			TrainingLiveStatusList: buildTrainingLiveStatusList(trainingIDs, statusSnapshotMap),
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
