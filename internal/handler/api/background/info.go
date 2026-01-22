package background

import (
	"honoka-chan/internal/schema/api/background"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func backgroundInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var backgroundList []int
	err = ss.MainEng.Table("background_m").Cols("background_id").Find(&backgroundList)
	if ss.CheckErr(err) {
		return
	}

	var backgroundID int
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Cols("background_id").Get(&backgroundID)
	if ss.CheckErr(err) {
		return
	}

	backgroundsList := []background.Info{}
	for _, id := range backgroundList {
		isSet := false
		if id == backgroundID {
			isSet = true
		}
		backgroundsList = append(backgroundsList, background.Info{
			BackgroundID: id,
			IsSet:        isSet,
			InsertDate:   time.Now().Format("2006-01-02 03:04:05"),
		})
	}

	res = background.InfoResp{
		Result: background.InfoData{
			BackgroundInfo: backgroundsList,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
