package user

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func UserApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "getNavi":
		res, err = userGetNavi(ctx)
	case "userInfo":
		res, err = userInfo(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("user", action)
	}
	return res, err
}
