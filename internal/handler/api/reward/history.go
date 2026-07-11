package reward

import (
	rewardapischema "honoka-chan/internal/schema/api/reward"
	"net/http"
	"time"
)

func rewardHistory() (res any, err error) {
	res = rewardapischema.HistoryResp{
		Result: rewardapischema.HistoryData{
			ItemCount: 0,
			Limit:     20,
			History:   []any{},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
