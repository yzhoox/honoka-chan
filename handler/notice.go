package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NoticeFriendVariety(ctx *gin.Context) {
	noticeResp := model.NoticeFriendVarietyResp{
		ResponseData: model.NoticeFriendVarietyRes{
			ItemCount:       1,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(noticeResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}

func NoticeFriendGreeting(ctx *gin.Context) {
	noticeResp := model.NoticeFriendGreetingResp{
		ResponseData: model.NoticeFriendGreetingRes{
			NextId:          0,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}
	resp, err := json.Marshal(noticeResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}

func NoticeUserGreeting(ctx *gin.Context) {
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
	resp, err := json.Marshal(noticeResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
