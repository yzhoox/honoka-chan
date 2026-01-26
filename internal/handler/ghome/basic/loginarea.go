package basic

import (
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func loginArea(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	ss.Respond(ghomeschema.LoginAreaResp{
		Code: 0,
		Msg:  "ok",
		Data: ghomeschema.LoginAreaData{
			UserID: ctx.PostForm("userid"),
		},
	})
}

func init() {
	router.AddHandler("v1", "POST", "/basic/loginarea", loginArea)
}
