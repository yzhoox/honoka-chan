package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func ProductList(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	prodReesp := model.ProductResp{
		ResponseData: model.ProductRes{
			RestrictionInfo: model.RestrictionInfo{
				Restricted: false,
			},
			UnderAgeInfo: model.UnderAgeInfo{
				BirthSet:    true,
				HasLimit:    false,
				LimitAmount: nil,
				MonthUsed:   0,
			},
			SnsProductList:   []any{},
			ProductList:      []any{},
			SubscriptionList: []any{},
			ShowPointShop:    true,
			ServerTimestamp:  time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(prodReesp)
}
