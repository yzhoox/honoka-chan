package payment

import (
	"honoka-chan/internal/schema/api/payment"
	"time"
)

func productList() (res any, err error) {
	res = payment.ProductListResp{
		Result: payment.ProductListData{
			RestrictionInfo: payment.RestrictionInfo{
				Restricted: false,
			},
			UnderAgeInfo: payment.UnderAgeInfo{
				BirthSet:    false,
				HasLimit:    false,
				LimitAmount: nil,
				MonthUsed:   0,
			},
			SnsProductList:   []payment.SnsProduct{},
			ProductList:      []payment.Product{},
			SubscriptionList: []payment.Subscription{},
			ShowPointShop:    false,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
