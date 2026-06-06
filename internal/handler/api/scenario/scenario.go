package scenario

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func ScenarioApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "scenarioStatus":
		res, err = scenarioStatus(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("scenario", action)
	}
	return res, err
}
