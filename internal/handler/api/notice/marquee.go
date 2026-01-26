package notice

import (
	noticeapischema "honoka-chan/internal/schema/api/notice"
	"time"
)

func noticeMarquee() (res any, err error) {
	res = noticeapischema.MarqueeResp{
		Result: noticeapischema.MarqueeData{
			ItemCount:   0,
			MarqueeList: []any{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
