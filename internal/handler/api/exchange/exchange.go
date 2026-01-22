package exchange

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ExchangeApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "owningPoint":
		res, err = owningPoint(ctx)
	default:
		err = fmt.Errorf("unimplemented action: exchange: %s", action)
	}
	return res, err
}
