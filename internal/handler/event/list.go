package event

import (
	"honoka-chan/internal/constant"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	commonschema "honoka-chan/internal/schema/common"
	eventschema "honoka-chan/internal/schema/event"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func list(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	targets := []eventschema.TargetList{}
	for i := range 6 {
		targets = append(targets, eventschema.TargetList{
			Position:      i + 1,
			IsDisplayable: false,
		})
	}

	ss.Respond(eventschema.ListResp{
		ResponseData: commonschema.ErrorData{
			ErrorCode: constant.ErrorCodeEventNoEventData,
		},
		ReleaseInfo: []any{},
		StatusCode:  600,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/event/eventList", middleware.Common, list)
}
