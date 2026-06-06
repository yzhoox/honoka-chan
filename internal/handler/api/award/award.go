package award

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func AwardApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "awardInfo":
		res, err = awardInfo(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("award", action)
	}
	return res, err
}
