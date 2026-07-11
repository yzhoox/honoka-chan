package banner

import (
	apibanner "honoka-chan/internal/handler/api/banner"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

type bannerResp struct {
	ResponseData any   `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}

func list(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	data, err := apibanner.BannerApi("bannerList")
	if ss.CheckErr(err) {
		return
	}

	resp, ok := data.(interface {
		GetResult() any
	})
	_ = resp
	_ = ok

	ss.Respond(bannerResp{
		ResponseData: extractBannerResult(data),
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func extractBannerResult(data any) any {
	if v, ok := data.(map[string]any); ok {
		if result, ok := v["result"]; ok {
			return result
		}
	}
	return extractResultField(data)
}

func init() {
	router.AddHandler("main.php", "POST", "/banner/bannerList", middleware.Common, list)
}
