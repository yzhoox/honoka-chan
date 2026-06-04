package profile

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func ProfileApi(ctx *gin.Context, action string, targetUserID int) (res any, err error) {
	switch action {
	case "cardRanking":
		res, err = cardRanking()
	case "liveCnt":
		res, err = liveCnt()
	case "profileInfo":
		res, err = profileInfo(ctx, targetUserID)
	default:
		err = fmt.Errorf("unimplemented action: profile: %s", action)
	}
	return res, err
}
