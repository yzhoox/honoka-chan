package download

import (
	"encoding/json"
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	downloadschema "honoka-chan/internal/schema/download"
	"honoka-chan/internal/session"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func update(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	downloadReq := downloadschema.UpdateReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []downloadschema.UpdateData{}
	if downloadReq.ExternalVersion != config.PackageVersion {
		pkgType := 99
		var pkgInfo []PkgInfo
		err := ss.MainEng.Table("download_m").Where("pkg_type = ? AND pkg_os = ?", pkgType, downloadReq.TargetOs).
			Cols("pkg_id,pkg_order,pkg_size").
			OrderBy("pkg_id ASC, pkg_order ASC").Find(&pkgInfo)
		if ss.CheckErr(err) {
			return
		}

		for _, pkg := range pkgInfo {
			pkgList = append(pkgList, downloadschema.UpdateData{
				Size: pkg.Size,
				URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
					config.Conf.Settings.CdnServer, downloadReq.TargetOs, pkgType, pkg.Id, pkg.Order),
				Version: config.PackageVersion,
			})
		}

		patchFileURL := fmt.Sprintf("%s/%s/archives/99_0_115.zip",
			config.Conf.Settings.CdnServer, downloadReq.TargetOs)
		resp, err := http.Get(patchFileURL)
		if err == nil {
			res, err := io.ReadAll(resp.Body)
			if err == nil {
				pkgList = append(pkgList, download.UpdateData{
					Size:    len(res),
					URL:     patchFileURL,
					Version: config.PackageVersion,
				})
			}
			defer resp.Body.Close()
		}
	}

	ss.Respond(downloadschema.UpdateResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/download/update", middleware.Common, update)
}
