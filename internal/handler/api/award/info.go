package award

import (
	awardapischema "honoka-chan/internal/schema/api/award"
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

	awardsList := []awardapischema.Info{}
	for _, id := range awardList {
		isSet := false
		if id == awardID {
			isSet = true
		}
		awardsList = append(awardsList, awardapischema.Info{
			AwardID:    id,
			IsSet:      isSet,
			InsertDate: "2023-03-20 03:58:55",
		})
	}

	res = awardapischema.InfoResp{
		Result: awardapischema.InfoData{
			AwardInfo: awardsList,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
