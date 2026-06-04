package payment

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func PaymentApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "productList":
		res, err = productList(ctx)
	default:
		err = fmt.Errorf("unimplemented action: payment: %s", action)
	}
	return res, err
}
