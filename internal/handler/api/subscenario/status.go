package subscenario

import (
	"honoka-chan/internal/schema/api/subscenario"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func subscenarioStatus(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var subScenarioID []int
	subScenarioLists := []subscenario.StatusList{}
	err = ss.MainEng.Table("subscenario_m").Cols("subscenario_id").OrderBy("subscenario_id ASC").Find(&subScenarioID)
	if ss.CheckErr(err) {
		return
	}

	for _, id := range subScenarioID {
		subScenarioLists = append(subScenarioLists, subscenario.StatusList{
			SubscenarioID: id,
			Status:        2,
		})
	}
	res = subscenario.StatusResp{
		Result: subscenario.StatusData{
			SubscenarioStatusList:  subScenarioLists,
			UnlockedSubscenarioIds: []any{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
