package exchange

import (
	exchangeapischema "honoka-chan/internal/schema/api/exchange"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func owningPoint(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var exchangeID []int
	exPointsList := []exchangeapischema.ExchangePointList{}
	err = ss.MainEng.Table("exchange_point_m").Cols("exchange_point_id").OrderBy("exchange_point_id ASC").Find(&exchangeID)
	if ss.CheckErr(err) {
		return
	}

	for _, id := range exchangeID {
		exPointsList = append(exPointsList, exchangeapischema.ExchangePointList{
			Rarity:        id,
			ExchangePoint: 9999,
		})
	}
	res = exchangeapischema.OwningPointResp{
		Result: exchangeapischema.OwningPointData{
			ExchangePointList: exPointsList,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
