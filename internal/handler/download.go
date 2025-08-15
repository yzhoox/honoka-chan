package handler

import (
	"encoding/json"
	"fmt"
	"honoka-chan/config"
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

type PkgInfo struct {
	Id    int `xorm:"pkg_id"`
	Order int `xorm:"pkg_order"`
	Size  int `xorm:"pkg_size"`
}

func DownloadAdditional(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	downloadReq := model.AdditionalReq{}
	err := json.Unmarshal([]byte(ctx.GetString("request_data")), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []model.AdditionalRes{}
	pkgType, pkgId := downloadReq.PackageType, downloadReq.PackageID
	var pkgInfo []PkgInfo
	err = ss.MainEng.Table("download_m").Where("pkg_type = ? AND pkg_id = ? AND pkg_os = ?", pkgType, pkgId, downloadReq.TargetOs).
		Cols("pkg_id,pkg_order,pkg_size").
		OrderBy("pkg_id ASC, pkg_order ASC").Find(&pkgInfo)
	if ss.CheckErr(err) {
		return
	}

	for _, pkg := range pkgInfo {
		pkgList = append(pkgList, model.AdditionalRes{
			Size: pkg.Size,
			URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
				config.Conf.Settings.CdnServer, downloadReq.TargetOs, pkgType, pkg.Id, pkg.Order),
		})
	}

	addResp := model.AdditionalResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(addResp)
}

func DownloadBatch(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	downloadReq := model.BatchReq{}
	err := json.Unmarshal([]byte(ctx.GetString("request_data")), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []model.BatchRes{}
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
			pkgList = append(pkgList, model.BatchRes{
				Size: pkg.Size,
				URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
					config.Conf.Settings.CdnServer, downloadReq.Os, pkgType, pkg.Id, pkg.Order),
			})
		}
	}

	batchResp := model.BatchResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(batchResp)
}

func DownloadUpdate(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	downloadReq := model.UpdateReq{}
	err := json.Unmarshal([]byte(ctx.GetString("request_data")), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	pkgList := []model.UpdateRes{}
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
			pkgList = append(pkgList, model.UpdateRes{
				Size: pkg.Size,
				URL: fmt.Sprintf("%s/%s/archives/%d_%d_%d.zip",
					config.Conf.Settings.CdnServer, downloadReq.TargetOs, pkgType, pkg.Id, pkg.Order),
				Version: config.PackageVersion,
			})
		}

		patchFileUrl := fmt.Sprintf("%s/%s/archives/99_0_115.zip",
			config.Conf.Settings.CdnServer, downloadReq.TargetOs)
		resp, err := http.Get(patchFileUrl)
		if err == nil {
			res, err := io.ReadAll(resp.Body)
			if err == nil {
				pkgList = append(pkgList, model.UpdateRes{
					Size:    len(res),
					URL:     patchFileUrl,
					Version: config.PackageVersion,
				})
			}
			defer resp.Body.Close()
		}
	}

	updateResp := model.UpdateResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(updateResp)
}

func DownloadUrl(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	downloadReq := model.UrlReq{}
	err := json.Unmarshal([]byte(ctx.GetString("request_data")), &downloadReq)
	if ss.CheckErr(err) {
		return
	}

	urlList := []string{}
	for _, v := range downloadReq.PathList {
		urlList = append(urlList, fmt.Sprintf("%s/%s/extracted/%s",
			config.Conf.Settings.BackupCdnServer, downloadReq.Os, strings.ReplaceAll(v, "\\", "")))
	}
	urlResp := model.UrlResp{
		ResponseData: model.UrlRes{
			UrlList: urlList,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(urlResp)
}

func DownloadEvent(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	eventResp := model.EventResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(eventResp)
}
