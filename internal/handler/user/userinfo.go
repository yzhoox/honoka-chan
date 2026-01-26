package user

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userInfo(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	pref := usermodel.UserPref{}
	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if ss.CheckErr(err) {
		return
	}

	userResp := userschema.UserInfoResp{
		ResponseData: userschema.UserInfoData{
			User: ss.GetUserInfo(),
			Birth: userschema.Birth{
				BirthMonth: 10,
				BirthDay:   18,
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
