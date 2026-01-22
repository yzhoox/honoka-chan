package liveicon

import (
	"honoka-chan/internal/schema/api/liveicon"
	"time"
)

func liveIconInfo() (res any, err error) {
	res = liveicon.InfoResp{
		Result: liveicon.InfoData{
			LiveNotesIconList: []int{1, 2, 3},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
