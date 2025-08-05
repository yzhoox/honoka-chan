package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"honoka-chan/internal/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func SetNotificationToken(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	notifResp := model.NotificationResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(notifResp)
}

func ChangeNavi(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	req := gjson.Parse(ctx.PostForm("request_data"))
	pref := tools.UserPref{
		UnitOwningUserID: int(req.Get("unit_owning_user_id").Int()),
	}
	_, err := ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Update(&pref)
	if ss.CheckErr(err) {
		return
	}

	naviResp := model.UserNaviChangeResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(naviResp)
}

func ChangeName(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	req := gjson.Parse(ctx.PostForm("request_data"))
	var oldName string
	exists, err := ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Cols("user_name").Get(&oldName)
	if ss.CheckErr(err) {
		return
	}
	if !exists {
		ctx.String(http.StatusForbidden, ErrorMsg)
		return
	}
	pref := tools.UserPref{
		UserName: req.Get("name").String(),
	}
	_, err = ss.UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Update(&pref)
	if ss.CheckErr(err) {
		return
	}
	nameResp := model.UserNameChangeResp{
		ResponseData: model.UserNameChangeRes{
			BeforeName:      oldName,
			AfterName:       req.Get("name").String(),
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(nameResp)
}
