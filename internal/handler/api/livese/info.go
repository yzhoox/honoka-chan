package livese

import (
	"honoka-chan/internal/schema/api/livese"
	"time"
)

func LiveSeInfo() (res any, err error) {
	res = livese.InfoResp{
		Result: livese.InfoData{
			LiveSeList: []int{1, 2, 3},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
