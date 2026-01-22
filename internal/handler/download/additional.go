package download

import (
	"encoding/json"
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/download"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

type PkgInfo struct {
	Id    int `xorm:"pkg_id"`
	Order int `xorm:"pkg_order"`
	Size  int `xorm:"pkg_size"`
}

func additional(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	downloadReq := download.AdditionalReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []download.AdditionalData{}
	pkgType, pkgId := downloadReq.PackageType, downloadReq.PackageID
	var pkgInfo []PkgInfo
	err = ss.MainEng.Table("download_m").Where("pkg_type = ? AND pkg_id = ? AND pkg_os = ?", pkgType, pkgId, downloadReq.TargetOs).
		Cols("pkg_id,pkg_order,pkg_size").
		OrderBy("pkg_id ASC, pkg_order ASC").Find(&pkgInfo)
	if ss.CheckErr(err) {
		return
	}

	for _, pkg := range pkgInfo {
		pkgList = append(pkgList, download.AdditionalData{
			Size: pkg.Size,
			URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
				config.Conf.Settings.CdnServer, downloadReq.TargetOs, pkgType, pkg.Id, pkg.Order),
		})
	}

	ss.Respond(download.AdditionalResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/download/additional", middleware.Common, additional)
}
