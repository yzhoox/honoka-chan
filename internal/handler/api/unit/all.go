package unit

import (
	"honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func unitAll(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	unitsData := []unit.Active{}
	err = ss.MainEng.Table("common_unit_m").Find(&unitsData)
	if ss.CheckErr(err) {
		return
	}

	userUnits := []unit.Active{}
	err = ss.UserEng.Table("user_unit").Where("user_id = ?", ss.UserID).Find(&userUnits)
	if ss.CheckErr(err) {
		return
	}
	unitsData = append(unitsData, userUnits...)

	res = unit.AllResp{
		Result: unit.AllData{
			Active:  unitsData,
			Waiting: []unit.Waiting{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
