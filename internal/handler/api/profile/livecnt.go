package profile

import (
	profileapischema "honoka-chan/internal/schema/api/profile"
	"time"
)

func liveCnt() (res any, err error) {
	res = profileapischema.LiveCntResp{
		Result: []profileapischema.LiveCntData{
			{
				Difficulty: 1,
				ClearCnt:   315,
			},
			{
				Difficulty: 2,
				ClearCnt:   310,
			},
			{
				Difficulty: 3,
				ClearCnt:   314,
			},
			{
				Difficulty: 4,
				ClearCnt:   455,
			},
			{
				Difficulty: 6,
				ClearCnt:   233,
			},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
