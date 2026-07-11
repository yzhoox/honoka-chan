package unit

import (
	unitapischema "honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func unitAccessoryAll(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	accessoryList := []unitapischema.AccessoryList{}
	err = ss.UserEng.Table("user_accessory").Where("user_id = ?", ss.UserID).Find(&accessoryList)
	if err != nil {
		return nil, err
	}
	for k := range accessoryList {
		accessoryList[k].NextExp = 0
		accessoryList[k].Level = 8
		accessoryList[k].MaxLevel = 8
		accessoryList[k].RankUpCount = 4
		accessoryList[k].FavoriteFlag = true
	}
	wearingInfo := []unitapischema.WearingInfo{}
	err = ss.UserEng.Table("user_accessory_wear").Where("user_id = ?", ss.UserID).Find(&wearingInfo)
	if err != nil {
		return nil, err
	}
	res = unitapischema.AccessoryAllResp{
		Result: unitapischema.AccessoryAllData{
			AccessoryList:      accessoryList,
			WearingInfo:        wearingInfo,
			EspecialCreateFlag: false,
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
