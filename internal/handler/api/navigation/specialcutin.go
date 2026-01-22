package navigation

import (
	"honoka-chan/internal/schema/api/navigation"
	"time"
)

func SpecialCutin() (res any, err error) {
	res = navigation.SpecialCutinResp{
		Result: navigation.SpecialCutinData{
			SpecialCutinList: []any{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
