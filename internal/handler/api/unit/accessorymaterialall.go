package unit

import (
	unitapischema "honoka-chan/internal/schema/api/unit"
	"net/http"
	"time"
)

func unitAccessoryMaterialAll() (res any, err error) {
	res = unitapischema.AccessoryMaterialAllResp{
		Result:     unitapischema.AccessoryMaterialAllData{},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
