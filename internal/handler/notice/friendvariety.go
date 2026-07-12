package notice

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	noticeschema "honoka-chan/internal/schema/notice"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func friendVariety(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	ss.Respond(noticeschema.FriendVarietyResp{
		ResponseData: noticeschema.FriendVarietyData{
			ItemCount:       1,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/notice/noticeFriendVariety", middleware.Common, friendVariety)
}
