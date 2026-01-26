package user

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func setNotificationToken(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(userschema.SetNotificationTokenResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/user/setNotificationToken", middleware.Common, setNotificationToken)
}
