package misc

import (
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func agreement(ctx *gin.Context) {
	ss := session.Attach(ctx)
	defer ss.FinalizeOrRollback()

	ss.Respond(ghomeschema.AgreementResp{
		ReturnCode:    0,
		ErrorType:     0,
		ReturnMessage: "",
		Data:          ghomeschema.AgreementData{},
	})
}

func init() {
	router.AddHandler("/", "GET", "/agreement/all", agreement)
}
