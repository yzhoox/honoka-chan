package unit

import (
	"honoka-chan/internal/schema/api/unit"
	"time"
)

func unitSupporterAll() (res any, err error) {
	res = unit.SupporterAllResp{
		Result: unit.SupporterAllData{
			UnitSupportList: []unit.SupporterList{},
		}, // 练习道具
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
