package notice

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	noticeschema "honoka-chan/internal/schema/notice"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func friendGreeting(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(noticeschema.FriendGreetingResp{
		ResponseData: noticeschema.FriendGreetingData{
			NextId:          0,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/notice/noticeFriendGreeting", middleware.Common, friendGreeting)
}
