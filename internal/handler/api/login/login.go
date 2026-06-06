package login

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func LoginApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "topInfo":
		res, err = loginTopInfo(ctx)
	case "topInfoOnce":
		res, err = loginTopInfoOnce(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("login", action)
	}
	return res, err
}
