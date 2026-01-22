package webui

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func logout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{
		Path:   "/admin",
		MaxAge: -1,
	})
	session.Save()

	ctx.Redirect(http.StatusFound, "/admin/login")
}

func init() {
	router.AddHandler("admin", "GET", "/logout", middleware.WebAuth, logout)
}
