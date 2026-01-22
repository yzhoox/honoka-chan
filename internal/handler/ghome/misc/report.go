package misc

import (
	"honoka-chan/internal/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func report(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func init() {
	router.AddHandler("/", "POST", "/report/ge/app", report)
}
