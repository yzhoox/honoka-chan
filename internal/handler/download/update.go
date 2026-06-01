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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const overrideServerConfigFileName = "99_0_115.zip"

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

		pkgMap := map[string]int{}
		for _, pkg := range pkgInfo {
			fileName := fmt.Sprintf("%d_%d_%d.zip", pkgType, pkg.Id, pkg.Order)
			url := fmt.Sprintf("%s/%s/archives/%s",
				config.Conf.Settings.CdnServer, downloadReq.TargetOs, fileName)

			pkgList = append(pkgList, downloadschema.UpdateData{
				Size:    pkg.Size,
				URL:     url,
				Version: config.PackageVersion,
			})
			pkgMap[fileName] = len(pkgList) - 1
		}

		serverConfigURL := fmt.Sprintf("%s/%s/archives/%s",
			config.Conf.Settings.CdnServer, downloadReq.TargetOs, overrideServerConfigFileName)
		serverConfigSize := getRemoteFileSize(serverConfigURL)

		overrideSource := getOverrideSource(downloadReq.TargetOs)
		if config.Conf.Settings.OverrideServerConfig.Enable && overrideSource.URL != "" {
			serverConfigURL = overrideSource.URL
			if overrideSource.Size > 0 {
				serverConfigSize = overrideSource.Size
			} else {
				overrideSize := getRemoteFileSize(serverConfigURL)
				if overrideSize > 0 {
					serverConfigSize = overrideSize
				}
			}
		}

		if index, ok := pkgMap[overrideServerConfigFileName]; ok {
			pkgList[index].URL = serverConfigURL
			if serverConfigSize > 0 {
				pkgList[index].Size = serverConfigSize
			}
		} else {
			pkgList = append(pkgList, downloadschema.UpdateData{
				Size:    serverConfigSize,
				URL:     serverConfigURL,
				Version: config.PackageVersion,
			})
		}

		applyOverrideFileSize(downloadReq.TargetOs, pkgList)
	}

	ss.Respond(downloadschema.UpdateResp{
		ResponseData: pkgList,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func getOverrideSource(targetOs string) config.OverrideFileSource {
	if strings.EqualFold(targetOs, "Android") {
		return config.Conf.Settings.OverrideServerConfig.Android
	}
	if strings.EqualFold(targetOs, "iOS") {
		return config.Conf.Settings.OverrideServerConfig.IOS
	}
	return config.OverrideFileSource{}
}

func getRemoteFileSize(url string) int {
	client := &http.Client{Timeout: 5 * time.Second}

	headReq, err := http.NewRequest(http.MethodHead, url, nil)
	if err == nil {
		if resp, err := client.Do(headReq); err == nil {
			_ = resp.Body.Close()
			if resp.ContentLength > 0 {
				return int(resp.ContentLength)
			}
		}
	}

	getResp, err := client.Get(url)
	if err != nil {
		return 0
	}
	defer getResp.Body.Close()

	dataLen, err := io.Copy(io.Discard, getResp.Body)
	if err != nil {
		return 0
	}

	return int(dataLen)
}

func applyOverrideFileSize(targetOs string, pkgList []downloadschema.UpdateData) {
	urlSizeCache := map[string]int{}

	for i := range pkgList {
		fileName := getFileNameFromURL(pkgList[i].URL)
		if fileName == "" {
			continue
		}

		for _, override := range config.Conf.Settings.OverrideFileSize {
			overrideOs := strings.TrimSpace(override.TargetOs)
			if overrideOs != "" && !strings.EqualFold(overrideOs, strings.TrimSpace(targetOs)) {
				continue
			}
			if strings.EqualFold(strings.TrimSpace(override.FileName), fileName) {
				size, ok := urlSizeCache[pkgList[i].URL]
				if !ok {
					size = getRemoteFileSize(pkgList[i].URL)
					urlSizeCache[pkgList[i].URL] = size
				}

				if size > 0 {
					pkgList[i].Size = size
				}
				break
			}
		}
	}
}

func getFileNameFromURL(rawURL string) string {
	urlNoQuery := strings.Split(rawURL, "?")[0]
	parts := strings.Split(urlNoQuery, "/")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[len(parts)-1])
}

func init() {
	router.AddHandler("main.php", "POST", "/download/update", middleware.Common, update)
}
