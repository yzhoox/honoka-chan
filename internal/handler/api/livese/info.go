package livese

import (
	liveseapischema "honoka-chan/internal/schema/api/livese"
	"net/http"
	"time"
)

func LiveSeInfo() (res any, err error) {
	res = liveseapischema.InfoResp{
		Result: liveseapischema.InfoData{
			LiveSeList: []int{1, 2, 3},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
