package payment

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func PaymentApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "productList":
		res, err = productList(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("payment", action)
	}
	return res, err
}
