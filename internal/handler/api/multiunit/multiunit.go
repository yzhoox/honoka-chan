package multiunit

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func MultiUnitApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "multiunitscenarioStatus":
		res, err = MultiUnitScenarioStatus(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("multiunit", action)
	}
	return res, err
}
