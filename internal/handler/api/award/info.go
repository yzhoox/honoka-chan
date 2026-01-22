package award

import (
	"honoka-chan/internal/schema/api/award"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func awardInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var awardList []int
	err = ss.MainEng.Table("award_m").Cols("award_id").Find(&awardList)
	if ss.CheckErr(err) {
		return
	}

	var awardID int
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Cols("award_id").Get(&awardID)
	if ss.CheckErr(err) {
		return
	}

	awardsList := []award.Info{}
	for _, id := range awardList {
		isSet := false
		if id == awardID {
			isSet = true
		}
		awardsList = append(awardsList, award.Info{
			AwardID:    id,
			IsSet:      isSet,
			InsertDate: time.Now().Format("2006-01-02 03:04:05"),
		})
	}

	res = award.InfoResp{
		Result: award.InfoData{
			AwardInfo: awardsList,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
