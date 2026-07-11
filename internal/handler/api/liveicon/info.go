package liveicon

import (
	liveiconapischema "honoka-chan/internal/schema/api/liveicon"
	"net/http"
	"time"
)

func liveIconInfo() (res any, err error) {
	res = liveiconapischema.InfoResp{
		Result: liveiconapischema.InfoData{
			LiveNotesIconList: []int{1, 2, 3},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
