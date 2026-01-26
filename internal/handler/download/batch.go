package download

import (
	"encoding/json"
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	downloadschema "honoka-chan/internal/schema/download"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func batch(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	downloadReq := downloadschema.BatchReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []downloadschema.BatchData{}
	if downloadReq.ClientVersion == config.PackageVersion {
		pkgType := downloadReq.PackageType
		var pkgInfo []PkgInfo
		err := ss.MainEng.Table("download_m").Where(builder.NotIn("pkg_id", downloadReq.ExcludedPackageIds)).Where("pkg_type = ? AND pkg_os = ?", pkgType, downloadReq.Os).
			Cols("pkg_id,pkg_order,pkg_size").
			OrderBy("pkg_id ASC, pkg_order ASC").Find(&pkgInfo)
		if ss.CheckErr(err) {
			return
		}

		for _, pkg := range pkgInfo {
			pkgList = append(pkgList, downloadschema.BatchData{
				Size: pkg.Size,
				URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
					config.Conf.Settings.CdnServer, downloadReq.Os, pkgType, pkg.Id, pkg.Order),
			})
		}
	}

	ss.Respond(downloadschema.BatchResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/download/batch", middleware.Common, batch)
}
