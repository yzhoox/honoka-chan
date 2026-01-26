package marathon

import (
	marathonapischema "honoka-chan/internal/schema/api/marathon"
	"time"
)

func marathonInfo() (res any, err error) {
	res = marathonapischema.InfoResp{
		Result:     []any{},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
