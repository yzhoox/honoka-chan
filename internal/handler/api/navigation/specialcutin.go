package navigation

import (
	navigationapischema "honoka-chan/internal/schema/api/navigation"
	"net/http"
	"time"
)

func SpecialCutin() (res any, err error) {
	res = navigationapischema.SpecialCutinResp{
		Result: navigationapischema.SpecialCutinData{
			SpecialCutinList: []any{},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
