package maintenance

import (
	"honoka-chan/internal/router"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func maintenance(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "common/maintenance.html", gin.H{
		"title":   "Love Live! 学园偶像祭",
		"content": template.HTML(`系统维护中，请稍后重试！`),
	})
}

func init() {
	router.AddHandler("/", "GET", "/resources/maintenace/maintenance.php", maintenance)
}
