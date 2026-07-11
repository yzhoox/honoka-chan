package marathon

import (
	marathonapischema "honoka-chan/internal/schema/api/marathon"
	"net/http"
	"time"
)

func marathonInfo() (res any, err error) {
	res = marathonapischema.InfoResp{
		Result:     []any{},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
