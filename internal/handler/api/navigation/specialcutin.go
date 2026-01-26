package navigation

import (
	navigationapischema "honoka-chan/internal/schema/api/navigation"
	"time"
)

func SpecialCutin() (res any, err error) {
	res = navigationapischema.SpecialCutinResp{
		Result: navigationapischema.SpecialCutinData{
			SpecialCutinList: []any{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
