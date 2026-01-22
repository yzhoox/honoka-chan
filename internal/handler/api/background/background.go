package background

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func BackgroundApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "backgroundInfo":
		res, err = backgroundInfo(ctx)
	default:
		err = fmt.Errorf("unimplemented action: background: %s", action)
	}
	return res, err
}
