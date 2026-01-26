package personalnotice

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	personalnoticeschema "honoka-chan/internal/schema/personalnotice"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func get(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	noticeResp := personalnoticeschema.GetResp{
		ResponseData: personalnoticeschema.GetData{
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

func init() {
	router.AddHandler("main.php", "POST", "/personalnotice/get", middleware.Common, get)
}
