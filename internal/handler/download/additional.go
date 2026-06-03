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
)

type PkgInfo struct {
	PkgType int `xorm:"pkg_type"`
	PkgID   int `xorm:"pkg_id"`
	Order   int `xorm:"pkg_order"`
	Size    int `xorm:"pkg_size"`
}

func additional(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	downloadReq := downloadschema.AdditionalReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &downloadReq)
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
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/download/additional", middleware.Common, additional)
}
