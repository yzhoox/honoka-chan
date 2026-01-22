package museum

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func MuseumApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "info":
		res, err = museumInfo(ctx)
	default:
		err = fmt.Errorf("unimplemented action: museum: %s", action)
	}
	return res, err
}
