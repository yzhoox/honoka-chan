package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func NoticeFriendVariety(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	noticeResp := model.NoticeFriendVarietyResp{
		ResponseData: model.NoticeFriendVarietyRes{
			ItemCount:       1,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(noticeResp)
}

func NoticeFriendGreeting(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	noticeResp := model.NoticeFriendGreetingResp{
		ResponseData: model.NoticeFriendGreetingRes{
			NextId:          0,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(noticeResp)
}

func NoticeUserGreeting(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	noticeResp := model.NoticeUserGreetingResp{
		ResponseData: model.NoticeUserGreetingRes{
			ItemCount:       0,
			HasNext:         false,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(noticeResp)
}
