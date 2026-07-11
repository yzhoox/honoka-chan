package subscenario

import (
	subscenarioapischema "honoka-chan/internal/schema/api/subscenario"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func subscenarioStatus(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var subScenarioID []int
	subScenarioLists := []subscenarioapischema.StatusList{}
	err = ss.MainEng.Table("subscenario_m").Cols("subscenario_id").OrderBy("subscenario_id ASC").Find(&subScenarioID)
	if err != nil {
		return nil, err
	}

	for _, id := range subScenarioID {
		subScenarioLists = append(subScenarioLists, subscenarioapischema.StatusList{
			SubscenarioID: id,
			Status:        2,
		})
	}
	res = subscenarioapischema.StatusResp{
		Result: subscenarioapischema.StatusData{
			SubscenarioStatusList:  subScenarioLists,
			UnlockedSubscenarioIds: []any{},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
