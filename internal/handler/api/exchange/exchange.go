package exchange

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func ExchangeApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "owningPoint":
		res, err = owningPoint(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("exchange", action)
	}
	return res, err
}
