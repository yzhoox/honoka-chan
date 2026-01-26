package user

import (
	usermodel "honoka-chan/internal/model/user"
	userapischema "honoka-chan/internal/schema/api/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	pref := usermodel.UserPref{}
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if err != nil {
		return nil, err
	}

	res = userapischema.InfoResp{
		Result:     ss.GetUserInfo(),
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
