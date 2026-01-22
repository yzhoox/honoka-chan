package event

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/event"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func list(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	targets := []event.TargetList{}
	for i := range 6 {
		targets = append(targets, event.TargetList{
			Position:      i + 1,
			IsDisplayable: false,
		})
	}

	// ss.Respond(event.ListResp{
	// 	ResponseData: map[string]int{
	// 		"error_code": 10004,
	// 	},
	// 	ReleaseInfo: []any{},
	// 	StatusCode:  600,
	// })

	ss.Respond(event.ListResp{
		ResponseData: event.ListData{
			TargetList:      targets,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/event/eventList", middleware.Common, list)
}
