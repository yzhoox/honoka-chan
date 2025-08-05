package handler

import (
	"encoding/json"
	"honoka-chan/internal/model"
	"honoka-chan/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func PersonalNotice(ctx *gin.Context) {
	noticeResp := model.PersonalNoticeResp{
		ResponseData: model.PersonalNoticeRes{
			HasNotice:       false,
			NoticeID:        0,
			Type:            0,
			Title:           "",
			Contents:        "",
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(noticeResp)
	utils.CheckErr(err)

	ctx.Header("X-Message-Sign", utils.GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
