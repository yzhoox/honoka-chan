package background

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func BackgroundApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "backgroundInfo":
		res, err = backgroundInfo(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("background", action)
	}
	return res, err
}
