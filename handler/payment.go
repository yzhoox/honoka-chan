package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ProductList(ctx *gin.Context) {
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
	resp, err := json.Marshal(prodReesp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
