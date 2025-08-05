package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"honoka-chan/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func ProfileRegister(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	req := gjson.Parse(ctx.PostForm("request_data"))
	pref := tools.UserPref{
		UserDesc: req.Get("introduction").String(),
	}

	_, err := ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Update(&pref)
	if ss.CheckErr(err) {
		return
	}

	profileResp := model.ProfileRegisterResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(profileResp)
}
