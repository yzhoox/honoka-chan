package subscenario

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func SubscenarioApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "subscenarioStatus":
		res, err = subscenarioStatus(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("subscenario", action)
	}
	return res, err
}
