package background

import (
	backgroundapischema "honoka-chan/internal/schema/api/background"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func backgroundInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var backgroundList []int
	err = ss.MainEng.Table("background_m").Cols("background_id").Find(&backgroundList)
	if err != nil {
		return nil, err
	}

	var backgroundID int
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Cols("background_id").Get(&backgroundID)
	if err != nil {
		return nil, err
	}

	backgroundsList := []backgroundapischema.Info{}
	for _, id := range backgroundList {
		isSet := false
		if id == backgroundID {
			isSet = true
		}
		backgroundsList = append(backgroundsList, backgroundapischema.Info{
			BackgroundID: id,
			IsSet:        isSet,
			InsertDate:   "2023-03-20 03:58:55",
		})
	}

	res = backgroundapischema.InfoResp{
		Result: backgroundapischema.InfoData{
			BackgroundInfo: backgroundsList,
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
