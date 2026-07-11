package achievement

import (
	achievementapischema "honoka-chan/internal/schema/api/achievement"
	"net/http"
	"time"
)

func unaccomplishList() (res any, err error) {
	listData := make([]achievementapischema.UnaccomplishListData, 0, len(achievementFilterCategoryList))
	for i := range achievementFilterCategoryList {
		listData = append(listData, achievementapischema.UnaccomplishListData{
			FilterCategoryID: achievementFilterCategoryList[i],
			AchievementList:  []achievementapischema.AchievementListItem{},
			Count:            0,
			IsLast:           true,
		})
	}

	res = achievementapischema.UnaccomplishListResp{
		Result:     listData,
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, nil
}
