package user

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userInfo(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	userResp := userschema.UserInfoResp{
		ResponseData: userschema.UserInfoData{
			User: ss.GetUserInfo(),
			Birth: userschema.Birth{
				BirthMonth: ss.UserPref.EffectiveBirthMonth(),
				BirthDay:   ss.UserPref.EffectiveBirthDay(),
			},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(userResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/user/userInfo", middleware.Common, userInfo)
}
