package unit

import (
	"fmt"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

type SetResp struct {
	ResponseData []any `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}

func favoriteAccessory(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	fmt.Println(ctx.MustGet("request_data").(string))

	ss.Respond(SetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/unit/favoriteAccessory", middleware.Common, favoriteAccessory)
}
