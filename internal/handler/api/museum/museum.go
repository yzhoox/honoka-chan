package museum

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func MuseumApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "info":
		res, err = museumInfo(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("museum", action)
	}
	return res, err
}
