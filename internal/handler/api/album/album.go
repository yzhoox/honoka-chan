package album

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func AlbumApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "albumAll":
		res, err = albumAll(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("album", action)
	}
	return res, err
}
