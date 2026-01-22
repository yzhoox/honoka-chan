package exchange

import (
	"honoka-chan/internal/schema/api/exchange"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func owningPoint(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var exchangeID []int
	exPointsList := []exchange.ExchangePointList{}
	err = ss.MainEng.Table("exchange_point_m").Cols("exchange_point_id").OrderBy("exchange_point_id ASC").Find(&exchangeID)
	if ss.CheckErr(err) {
		return
	}

	for _, id := range exchangeID {
		exPointsList = append(exPointsList, exchange.ExchangePointList{
			Rarity:        id,
			ExchangePoint: 9999,
		})
	}
	res = exchange.OwningPointResp{
		Result: exchange.OwningPointData{
			ExchangePointList: exPointsList,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
