package achievement

import (
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

var achievementFilterCategoryList = []int{
	1, // 重要
	3, // 挑战
	4, // 每日
	5, // 期间限定
	6, // 8周年限定/全世界5000万人突破纪念
	7, // 回忆画廊
	8, // 9周年限定
	9, // 连携课题
}

func AchievementApi(ctx *gin.Context, action string) (res any, err error) {
	switch action {
	case "unaccomplishList":
		res, err = unaccomplishList()
	case "initialAccomplishedList":
		res, err = initialAccomplishedList(ctx)
	default:
		err = honokautils.NewUnimplementedActionError("achievement", action)
	}
	return res, err
}
