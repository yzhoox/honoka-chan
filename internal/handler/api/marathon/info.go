package marathon

import (
	"honoka-chan/internal/schema/api/marathon"
	"time"
)

func marathonInfo() (res any, err error) {
	res = marathon.InfoResp{
		Result:     []any{},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
