package reward

import (
	rewardapischema "honoka-chan/internal/schema/api/reward"
	"time"
)

func rewardHistory() (res any, err error) {
	res = rewardapischema.HistoryResp{
		Result: rewardapischema.HistoryData{
			ItemCount: 0,
			Limit:     20,
			History:   []any{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
