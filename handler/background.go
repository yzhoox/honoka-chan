package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"honoka-chan/tools"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func BackgroundSet(ctx *gin.Context) {
	req := gjson.Parse(ctx.PostForm("request_data"))
	pref := tools.UserPref{
		BackgroundID: int(req.Get("background_id").Int()),
	}
	_, err := UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Update(&pref)
	CheckErr(err)
	backgroundResp := model.BackgroundSetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}
	resp, err := json.Marshal(backgroundResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
