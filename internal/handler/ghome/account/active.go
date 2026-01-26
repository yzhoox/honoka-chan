package account

import (
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func active(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	ss.Respond(ghomeschema.ActiveResp{
		Code: 0,
		Msg:  "ok",
		Data: ghomeschema.ActiveData{
			Message: "ok",
			Result:  0,
		},
	})
}

func init() {
	router.AddHandler("v1", "POST", "/account/active", active)
}
