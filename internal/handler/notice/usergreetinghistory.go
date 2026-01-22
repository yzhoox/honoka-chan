package notice

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/notice"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userGreetingHistory(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	ss.Respond(notice.UserGreetingResp{
		ResponseData: notice.UserGreetingData{
			ItemCount:       0,
			HasNext:         false,
			NoticeList:      []any{},
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/notice/noticeUserGreetingHistory", middleware.Common, userGreetingHistory)
}
