package unit

import (
	unitapischema "honoka-chan/internal/schema/api/unit"
	"net/http"
	"time"
)

func unitSupporterAll() (res any, err error) {
	res = unitapischema.SupporterAllResp{
		Result: unitapischema.SupporterAllData{
			UnitSupportList: []unitapischema.SupporterList{},
		}, // 练习道具
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
