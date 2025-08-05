package handler

import (
	"encoding/json"
	"honoka-chan/internal/model"
	"honoka-chan/internal/tools"
	"honoka-chan/internal/utils"
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
	utils.CheckErr(err)
	backgroundResp := model.BackgroundSetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}
	resp, err := json.Marshal(backgroundResp)
	utils.CheckErr(err)

	ctx.Header("X-Message-Sign", utils.GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
