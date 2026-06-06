package announce

import (
	"honoka-chan/internal/router"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

const announceContent = `目前开发完毕的主要功能包括：<br><ul><li>登录与基础用户数据</li><li>相册、编队、饰品、宝石</li><li>Live 相关主要流程</li><li>个人信息设置</li><li>好友系统</li><li>招募系统</li><li>官方漫画</li></ul><br>
			其他功能仍在开发中，有报错属于正常现象。<br><br>
			侵权联系：<a href="https://space.bilibili.com/671443">梦路_YumeMichi @bilibili</a> 进行删除。`

func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "common/announce.html", gin.H{
		"title":   "Love Live! 学园偶像祭",
		"content": template.HTML(announceContent),
	})
}

func init() {
	router.AddHandler("webview.php", "GET", "/announce/index", index)
}
