package misc

import (
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func appReport(ctx *gin.Context) {
	ss := session.Attach(ctx)
	defer ss.Finalize()

	ss.Respond(ghomeschema.AppReportResp{
		Code: 0,
		Msg:  "",
		Data: ghomeschema.AppReportData{
			NeedReport: 0,
		},
	})
}

func init() {
	router.AddHandler("/", "GET", "/integration/appReport/initialize", appReport)
}
