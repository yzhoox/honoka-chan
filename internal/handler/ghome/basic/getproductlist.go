package basic

import (
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func getProductList(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	getProductListData := ghome.GetProductListData{
		Message: []string{},
		Result:  0,
	}

	ss.Respond(ghome.GetProductListResp{
		Code: 1,
		Msg:  "ok",
		Data: getProductListData,
	})
}

func init() {
	router.AddHandler("v1", "POST", "/basic/getProductList", getProductList)
}
