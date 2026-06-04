package login

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LoginApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "topInfo":
		res, err = loginTopInfo(ctx)
	case "topInfoOnce":
		res, err = loginTopInfoOnce()
	default:
		err = fmt.Errorf("unimplemented action: login: %s", action)
	}
	return res, err
}
