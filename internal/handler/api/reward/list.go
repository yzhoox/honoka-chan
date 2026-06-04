package reward

import (
	rewardapischema "honoka-chan/internal/schema/api/reward"
	"time"
)

func rewardList() (res any, err error) {
	res = rewardapischema.ListResp{
		Result: rewardapischema.ListData{
			ItemCount: 0,
			Limit:     20,
			Order:     0,
			Items:     []any{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
