package download

import (
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	downloadschema "honoka-chan/internal/schema/download"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const overrideServerConfigFileName = "99_0_115.zip"

func update(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	downloadReq := downloadschema.UpdateReq{}
	err := honokautils.ParseRequestData(ctx, &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []downloadschema.UpdateData{}
	if downloadReq.ExternalVersion != config.PackageVersion {
		pkgType := 99
		var pkgInfo []PkgInfo
		err := ss.MainEng.Table("download_m").Where("pkg_type = ? AND pkg_os = ?", pkgType, downloadReq.TargetOs).
			OrderBy("pkg_id ASC, pkg_order ASC").Find(&pkgInfo)
		if ss.CheckErr(err) {
			return
		}

		for _, pkg := range pkgInfo {
			fileName := fmt.Sprintf("%d_%d_%d.zip", pkg.PkgType, pkg.PkgID, pkg.Order)
			url := fmt.Sprintf("%s/%s/archives/%s",
				config.Conf.Settings.CdnServer, downloadReq.TargetOs, fileName)

			pkgList = append(pkgList, downloadschema.UpdateData{
				Size:    pkg.Size,
				URL:     url,
				Version: config.PackageVersion,
			})
		}

		serverConfigURL := fmt.Sprintf("%s/%s/archives/%s",
			config.Conf.Settings.CdnServer, downloadReq.TargetOs, overrideServerConfigFileName)
		serverConfigSize, ok := getRemoteFileSize(serverConfigURL)
		if ok {
			pkgList = append(pkgList, downloadschema.UpdateData{
				Size:    serverConfigSize,
				URL:     serverConfigURL,
				Version: config.PackageVersion,
			})
		}
	}

	ss.Respond(downloadschema.UpdateResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func getRemoteFileSize(url string) (int, bool) {
	client := &http.Client{Timeout: 5 * time.Second}

	headReq, err := http.NewRequest(http.MethodHead, url, nil)
	if err == nil {
		if resp, err := client.Do(headReq); err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= http.StatusOK &&
				resp.StatusCode < http.StatusMultipleChoices && resp.ContentLength > 0 {
				return int(resp.ContentLength), true
			}
		}
	}

	getResp, err := client.Get(url)
	if err != nil {
		return 0, false
	}
	defer getResp.Body.Close()
	if getResp.StatusCode < http.StatusOK || getResp.StatusCode >= http.StatusMultipleChoices {
		return 0, false
	}

	dataLen, err := io.Copy(io.Discard, getResp.Body)
	if err != nil {
		return 0, false
	}

	return int(dataLen), dataLen > 0
}

func init() {
	router.AddHandler("main.php", "POST", "/download/update", middleware.Common, update)
}
