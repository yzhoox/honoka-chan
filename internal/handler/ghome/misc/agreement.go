package misc

import (
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func agreement(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	ss.Respond(ghome.AgreementResp{
		ReturnCode:    0,
		ErrorType:     0,
		ReturnMessage: "",
		Data:          ghome.AgreementData{},
	})
}

func init() {
	router.AddHandler("/", "GET", "/agreement/all", agreement)
}
