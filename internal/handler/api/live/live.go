package live

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func LiveApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "liveStatus":
		res, err = liveStatus(ctx)
	case "schedule":
		res, err = liveSchedule(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("live", action)
	}
	return res, err
}
