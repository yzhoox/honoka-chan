package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func UserApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "getNavi":
		res, err = userGetNavi(ctx)
	case "userInfo":
		res, err = userInfo(ctx)
	default:
		err = fmt.Errorf("unimplemented action: user: %s", action)
	}
	return res, err
}
