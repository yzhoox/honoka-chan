package award

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AwardApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "awardInfo":
		res, err = awardInfo(ctx)
	default:
		err = fmt.Errorf("unimplemented action: award: %s", action)
	}
	return res, err
}
