package payment

import (
	paymentapischema "honoka-chan/internal/schema/api/payment"
	"time"
)

func productList() (res any, err error) {
	res = paymentapischema.ProductListResp{
		Result: paymentapischema.ProductListData{
			RestrictionInfo: paymentapischema.RestrictionInfo{
				Restricted: false,
			},
			UnderAgeInfo: paymentapischema.UnderAgeInfo{
				BirthSet:    false,
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
