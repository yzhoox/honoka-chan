package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthKey(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	authResp := model.AuthKeyResp{
		ResponseData: model.AuthKeyRes{
			AuthorizeToken: ctx.GetString("authorize_token"),
			DummyToken:     ctx.GetString("dummy_token"),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(authResp)
}

func Login(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	loginKey := ctx.GetString("login_key")
	var userId int
	exists, err := ss.UserEng.Table("user_key").Where("key = ?", loginKey).Cols("userid").Get(&userId)
	if ss.CheckErr(err) {
		return
	}

	if !exists || userId == 0 {
		userId = 9999999
	}
	ctx.Set("userid", userId)

	err = db.DB.Set([]byte(strconv.Itoa(userId)), []byte(ctx.GetString("authorize_token")))
	if ss.CheckErr(err) {
		return
	}

	loginResp := model.LoginResp{
		ResponseData: model.LoginRes{
			AuthorizeToken:  ctx.GetString("authorize_token"),
			UserId:          userId,
			ServerTimestamp: time.Now().Unix(),
			AdultFlag:       2,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(loginResp)
}
