package user

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setNotificationToken(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	ss.Respond(userschema.SetNotificationTokenResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/user/setNotificationToken", middleware.Common, setNotificationToken)
}
