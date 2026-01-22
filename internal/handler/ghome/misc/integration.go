package misc

import (
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func appReport(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	ss.Respond(ghome.AppReportResp{
		Code: 0,
		Msg:  "",
		Data: ghome.AppReportData{
			NeedReport: 0,
		},
	})
}

func init() {
	router.AddHandler("/", "GET", "/integration/appReport/initialize", appReport)
}
