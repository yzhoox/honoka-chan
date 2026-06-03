package guest

import (
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func status(ctx *gin.Context) {
	ss := session.Attach(ctx)
	defer ss.Finalize()

	ss.Respond(ghomeschema.GuestStatusResp{
		Code: 0,
		Msg:  "ok",
		Data: ghomeschema.GuestStatusData{
			Disablead:   1,
			Loginswitch: 1,
			Message:     "ok",
			Result:      0,
		},
	})
}

func init() {
	router.AddHandler("v1", "POST", "/guest/status", status)
}
