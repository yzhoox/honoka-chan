package user

import (
	userapischema "honoka-chan/internal/schema/api/user"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	res = userapischema.InfoResp{
		Result: userschema.UserInfoData{
			User: ss.GetUserInfo(),
			Birth: userschema.Birth{
				BirthMonth: ss.UserPref.EffectiveBirthMonth(),
				BirthDay:   ss.UserPref.EffectiveBirthDay(),
			},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
