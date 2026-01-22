package webui

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/webui"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
)

func login(ctx *gin.Context) {
	area := ctx.PostForm("area")
	user := ctx.PostForm("user")
	pass := ctx.PostForm("pass")
	if area == "" || user == "" || pass == "" {
		ctx.JSON(http.StatusOK, webui.Msg{
			Code:     1,
			Message:  "参数不完整！",
			Redirect: "",
		})
		return
	}

	userName := " " + area + "-" + user
	var userId int
	exists, err := db.UserEng.Table("users").Where("phone = ? AND password = ?", userName, openssl.Md5ToString(pass)).Cols("user_id").Get(&userId)
	utils.CheckErr(err)
	if !exists {
		ctx.JSON(http.StatusOK, webui.Msg{
			Code:     1,
			Message:  "账号不存在或者密码有误！",
			Redirect: "",
		})
		return
	}

	session := sessions.Default(ctx)
	session.Options(sessions.Options{
		MaxAge: 3600 * 24,
	})
	session.Set("userid", userId)
	session.Save()

	ctx.JSON(http.StatusOK, webui.Msg{
		Code:     0,
		Message:  "登录成功！",
		Redirect: "/admin/index",
	})
}

func init() {
	router.AddHandler("admin", "GET", "/login", middleware.WebAuth, func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "admin/login.html", gin.H{})
	})
	router.AddHandler("admin", "POST", "/login", middleware.WebAuth, login)
}
