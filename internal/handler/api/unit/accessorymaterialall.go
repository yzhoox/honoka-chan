package unit

import (
	unitapischema "honoka-chan/internal/schema/api/unit"
	"time"
)

func unitAccessoryMaterialAll() (res any, err error) {
	res = unitapischema.AccessoryMaterialAllResp{
		Result:     unitapischema.AccessoryMaterialAllData{},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
