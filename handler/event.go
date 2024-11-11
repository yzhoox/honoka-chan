package handler

import (
	"encoding/json"
	"honoka-chan/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func EventList(ctx *gin.Context) {
	targets := []model.TargetList{}
	for i := 0; i < 6; i++ {
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
	resp, err := json.Marshal(eventsResp)
	CheckErr(err)

	ctx.Header("X-Message-Sign", GenXMS(resp))
	ctx.String(http.StatusOK, string(resp))
}
