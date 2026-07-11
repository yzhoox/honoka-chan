package download

import (
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	downloadschema "honoka-chan/internal/schema/download"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func additional(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	downloadReq := downloadschema.AdditionalReq{}
	err := honokautils.ParseRequestData(ctx, &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []downloadschema.AdditionalData{}
	pkgType, pkgId := downloadReq.PackageType, downloadReq.PackageID
	var pkgInfo []PkgInfo
	err = ss.MainEng.Table("download_m").Where("pkg_type = ? AND pkg_id = ? AND pkg_os = ?", pkgType, pkgId, downloadReq.TargetOs).
		OrderBy("pkg_id ASC, pkg_order ASC").Find(&pkgInfo)
	if ss.CheckErr(err) {
		return
	}

	for _, pkg := range pkgInfo {
		pkgList = append(pkgList, downloadschema.AdditionalData{
			Size: pkg.Size,
			URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
				config.Conf.Settings.CdnServer, downloadReq.TargetOs, pkg.PkgType, pkg.PkgID, pkg.Order),
		})
	}

	ss.Respond(downloadschema.AdditionalResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/download/additional", middleware.Common, additional)
}
