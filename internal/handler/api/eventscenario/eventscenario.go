package eventscenario

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func EventScenarioApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "status":
		res, err = eventScenarioStatus(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("eventscenario", action)
	}
	return res, err
}
