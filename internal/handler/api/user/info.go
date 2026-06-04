package user

import (
	userapischema "honoka-chan/internal/schema/api/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	res = userapischema.InfoResp{
		Result:     ss.GetUserInfo(),
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
