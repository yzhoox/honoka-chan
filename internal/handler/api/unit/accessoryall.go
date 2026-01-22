package unit

import (
	"honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func unitAccessoryAll(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	accessoryList := []unit.AccessoryList{}
	err = ss.MainEng.Table("common_accessory_m").Find(&accessoryList)
	if ss.CheckErr(err) {
		return
	}
	for k := range accessoryList {
		accessoryList[k].NextExp = 0
		accessoryList[k].Level = 8
		accessoryList[k].MaxLevel = 8
		accessoryList[k].RankUpCount = 4
		accessoryList[k].FavoriteFlag = true
	}
	wearingInfo := []unit.WearingInfo{}
	err = ss.UserEng.Table("user_accessory_wear").Where("user_id = ?", ss.UserID).Find(&wearingInfo)
	if ss.CheckErr(err) {
		return
	}
	res = unit.AccessoryAllResp{
		Result: unit.AccessoryAllData{
			AccessoryList:      accessoryList,
			WearingInfo:        wearingInfo,
			EspecialCreateFlag: false,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
