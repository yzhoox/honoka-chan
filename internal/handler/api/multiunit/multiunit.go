package multiunit

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func MultiUnitApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "multiunitscenarioStatus":
		res, err = MultiUnitScenarioStatus(ctx)
	default:
		err = fmt.Errorf("unimplemented action: multiunit: %s", action)
	}
	return res, err
}
