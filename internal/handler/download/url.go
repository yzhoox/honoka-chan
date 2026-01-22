package download

import (
	"encoding/json"
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/download"
	"honoka-chan/internal/session"
	"strings"

	"github.com/gin-gonic/gin"
)

func url(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	downloadReq := download.UrlReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	cdnServer := config.Conf.Settings.BackupCdnServer
	if cdnServer == "" {
		cdnServer = config.Conf.Settings.CdnServer
	}

	urlList := []string{}
	for _, v := range downloadReq.PathList {
		urlList = append(urlList, fmt.Sprintf("%s/%s/extracted/%s",
			cdnServer, downloadReq.Os, strings.ReplaceAll(v, "\\", "")))
	}

	ss.Respond(download.UrlResp{
		ResponseData: download.UrlData{
			UrlList: urlList,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/download/url", middleware.Common, url)
}
