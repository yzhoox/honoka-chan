package payment

import (
	paymentapischema "honoka-chan/internal/schema/api/payment"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func productList(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	res = paymentapischema.ProductListResp{
		Result: paymentapischema.ProductListData{
			RestrictionInfo: paymentapischema.RestrictionInfo{
				Restricted: false,
			},
			UnderAgeInfo: paymentapischema.UnderAgeInfo{
				BirthSet:    ss.UserPref.HasBirthDate(),
				HasLimit:    false,
				LimitAmount: nil,
				MonthUsed:   0,
			},
			SnsProductList:   []paymentapischema.SnsProduct{},
			ProductList:      []paymentapischema.Product{},
			SubscriptionList: []paymentapischema.Subscription{},
			ShowPointShop:    false,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
