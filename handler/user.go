package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"honoka-chan/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func SetNotificationToken(ctx *gin.Context) {
	notifResp := model.NotificationResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}
	resp, err := json.Marshal(notifResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}

func ChangeNavi(ctx *gin.Context) {
	req := gjson.Parse(ctx.PostForm("request_data"))
	pref := tools.UserPref{
		UnitOwningUserID: int(req.Get("unit_owning_user_id").Int()),
	}
	_, err := UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Update(&pref)
	CheckErr(err)
	naviResp := model.UserNaviChangeResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}
	resp, err := json.Marshal(naviResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}

func ChangeName(ctx *gin.Context) {
	req := gjson.Parse(ctx.PostForm("request_data"))
	var oldName string
	exists, err := UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Cols("user_name").Get(&oldName)
	CheckErr(err)
	if !exists {
		ctx.String(http.StatusForbidden, ErrorMsg)
		return
	}
	pref := tools.UserPref{
		UserName: req.Get("name").String(),
	}
	_, err = UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Update(&pref)
	CheckErr(err)
	nameResp := model.UserNameChangeResp{
		ResponseData: model.UserNameChangeRes{
			BeforeName:      oldName,
			AfterName:       req.Get("name").String(),
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(nameResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
