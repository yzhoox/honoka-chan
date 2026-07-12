package basic

import (
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func getProductList(ctx *gin.Context) {
	ss := session.Attach(ctx)
	defer ss.FinalizeOrRollback()

	getProductListData := ghomeschema.GetProductListData{
		Message: []string{},
		Result:  0,
	}

	ss.Respond(ghomeschema.GetProductListResp{
		Code: 1,
		Msg:  "ok",
		Data: getProductListData,
	})
}

func init() {
	router.AddHandler("v1", "POST", "/basic/getProductList", getProductList)
}
