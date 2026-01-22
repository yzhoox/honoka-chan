package live

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LiveApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "liveStatus":
		res, err = liveStatus(ctx)
	case "schedule":
		res, err = liveSchedule(ctx)
	default:
		err = fmt.Errorf("unimplemented action: live: %s", action)
	}
	return res, err
}
