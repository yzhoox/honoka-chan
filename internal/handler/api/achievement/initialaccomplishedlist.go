package achievement

import (
	achievementapischema "honoka-chan/internal/schema/api/achievement"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func initialAccomplishedList(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	grouped, err := ListVisibleAccomplishedAchievements(ss)
	if err != nil {
		return nil, err
	}

	listData := make([]achievementapischema.UnaccomplishListData, 0, len(achievementFilterCategoryList))
	for _, filterCategoryID := range achievementFilterCategoryList {
		achievementItems := grouped[filterCategoryID]
		if achievementItems == nil {
			achievementItems = []achievementapischema.AchievementListItem{}
		} else {
			if len(achievementItems) > 5 {
				achievementItems = achievementItems[:5]
			}
		}

		listData = append(listData, achievementapischema.UnaccomplishListData{
			FilterCategoryID: filterCategoryID,
			AchievementList:  achievementItems,
			Count:            len(grouped[filterCategoryID]),
			IsLast:           true,
		})
	}

	return achievementapischema.UnaccomplishListResp{
		Result:     listData,
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}, nil
}
