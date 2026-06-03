package scenario

import (
	scenarioapischema "honoka-chan/internal/schema/api/scenario"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func scenarioStatus(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var scenarioID []int
	scenarioLists := []scenarioapischema.StatusList{}
	err = ss.MainEng.Table("scenario_m").Cols("scenario_id").OrderBy("scenario_id ASC").Find(&scenarioID)
	if err != nil {
		return nil, err
	}

	for _, id := range scenarioID {
		scenarioLists = append(scenarioLists, scenarioapischema.StatusList{
			ScenarioID: id,
			Status:     2,
		})
	}
	res = scenarioapischema.StatusResp{
		Result: scenarioapischema.StatusData{
			ScenarioStatusList: scenarioLists,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
