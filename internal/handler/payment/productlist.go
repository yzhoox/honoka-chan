package payment

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	paymentapischema "honoka-chan/internal/schema/api/payment"
	"honoka-chan/internal/schema/payment"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func productList(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	prodResp := payment.ProductListResp{
		ResponseData: payment.ProductListData{
			RestrictionInfo: paymentapischema.RestrictionInfo{
				Restricted: false,
			},
			UnderAgeInfo: paymentapischema.UnderAgeInfo{
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

	ss.Respond(prodResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/payment/productList", middleware.Common, productList)
}
