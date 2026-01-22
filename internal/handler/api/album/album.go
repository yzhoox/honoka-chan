package album

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AlbumApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "albumAll":
		res, err = albumAll(ctx)
	default:
		err = fmt.Errorf("unimplemented action: album: %s", action)
	}
	return res, err
}
