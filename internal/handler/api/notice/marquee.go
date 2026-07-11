package notice

import (
	noticeapischema "honoka-chan/internal/schema/api/notice"
	"net/http"
	"time"
)

func noticeMarquee() (res any, err error) {
	res = noticeapischema.MarqueeResp{
		Result: noticeapischema.MarqueeData{
			ItemCount:   0,
			MarqueeList: []any{},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
