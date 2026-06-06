package profile

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func ProfileApi(ctx *gin.Context, action string, targetUserID int) (res any, err error) {
	switch action {
	case "cardRanking":
		res, err = cardRanking()
	case "liveCnt":
		res, err = liveCnt(ctx, targetUserID)
	case "profileInfo":
		res, err = profileInfo(ctx, targetUserID)
	default:
		err = honokautils.NewUnimplementedActionError("profile", action)
	}
	return res, err
}
