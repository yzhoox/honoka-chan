package download

import (
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	downloadschema "honoka-chan/internal/schema/download"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func batch(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	downloadReq := downloadschema.BatchReq{}
	err := honokautils.ParseRequestData(ctx, &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []downloadschema.BatchData{}
	if downloadReq.ClientVersion == config.PackageVersion {
		pkgType := downloadReq.PackageType

		// 首次下载时客户端发送的 downloadReq.PackageType = 0 但是实际下载的是 [0,4]
		// 此时直接使用 downloadReq.PackageType 会导致下载数据量计算错误
		// 游戏会自动将 [0,4] 的包大小合并计算
		// downloadReq.PackageType = 0 的包已经在 /download/additional 中获取
		// 所以这里只需要返回 downloadReq.PackageType = 4 的包即可
		if downloadReq.PackageType == 0 {
			pkgType = 4
		}

		var pkgInfo []PkgInfo
		err := ss.MainEng.Table("download_m").
			Where(builder.NotIn("pkg_id", downloadReq.ExcludedPackageIds)).
			Where("pkg_type = ? AND pkg_os = ?", pkgType, downloadReq.Os).
			OrderBy("pkg_id ASC, pkg_order ASC").Find(&pkgInfo)
		if ss.CheckErr(err) {
			return
		}

		for _, pkg := range pkgInfo {
			pkgList = append(pkgList, downloadschema.BatchData{
				Size: pkg.Size,
				URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
					config.Conf.Settings.CdnServer, downloadReq.Os, pkg.PkgType, pkg.PkgID, pkg.Order),
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
