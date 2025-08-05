package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func EventList(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	targets := []model.TargetList{}
	for i := range 6 {
		targets = append(targets, model.TargetList{
			Position:      i + 1,
			IsDisplayable: false,
		})
	}
	eventsResp := model.EventsResp{
		ResponseData: model.EventsRes{
			TargetList:      targets,
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(eventsResp)
}
