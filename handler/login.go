package handler

import (
	"encoding/json"
	"honoka-chan/db"
	"honoka-chan/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthKey(ctx *gin.Context) {
	authResp := model.AuthKeyResp{
		ResponseData: model.AuthKeyRes{
			AuthorizeToken: ctx.GetString("authorize_token"),
			DummyToken:     ctx.GetString("dummy_token"),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(authResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.JSON(http.StatusOK, authResp)
}

func Login(ctx *gin.Context) {
	loginKey := ctx.GetString("login_key")
	var userId int
	exists, err := UserEng.Table("user_key").Where("key = ?", loginKey).Cols("userid").Get(&userId)
	CheckErr(err)

	if !exists || userId == 0 {
		userId = 9999999
	}
	ctx.Set("userid", userId)

	err = db.DB.Set([]byte(strconv.Itoa(userId)), []byte(ctx.GetString("authorize_token")))
	CheckErr(err)

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
	resp, err := json.Marshal(loginResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.JSON(http.StatusOK, loginResp)
}
