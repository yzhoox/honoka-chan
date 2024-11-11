package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"honoka-chan/tools"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func AwardSet(ctx *gin.Context) {
	req := gjson.Parse(ctx.PostForm("request_data"))
	pref := tools.UserPref{
		AwardID: int(req.Get("award_id").Int()),
	}
	_, err := UserEng.Table("user_preference_m").Where("user_id = ?", ctx.GetString("userid")).Update(&pref)
	CheckErr(err)
	awardResp := model.AwardSetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}
	resp, err := json.Marshal(awardResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
