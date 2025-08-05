package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func PersonalNotice(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

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

	ss.Respond(noticeResp)
}
