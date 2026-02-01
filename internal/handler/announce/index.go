package announce

import (
	"honoka-chan/internal/router"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "common/announce.html", gin.H{
		"title": "Love Live! 学园偶像祭",
		"content": template.HTML(`目前开发完毕的功能包括：<br><ul><li>登录</li><li>相册</li><li>编队</li><li>饰品</li><li>宝石</li><li>Live</li><li>个人信息设置</li><li>官方漫画</li><ul><br>
			其他功能仍在开发中，有报错属于正常现象。<br><br>
			侵权联系：<a href="https://space.bilibili.com/671443">梦路_YumeMichi @bilibili</a> 进行删除。`),
	})
}

func init() {
	router.AddHandler("webview.php", "GET", "/announce/index", index)
}
