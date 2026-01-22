package webui

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func card(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/card.html", gin.H{
		"menu": 1,
		"url":  strings.Split(ctx.Request.URL.String(), "?")[0],
	})
}

func index(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "admin/index.html", gin.H{
		"url": strings.Split(ctx.Request.URL.String(), "?")[0],
	})
}

func init() {
	router.AddHandler("admin", "GET", "/card", middleware.WebAuth, card)
	router.AddHandler("admin", "GET", "/index", middleware.WebAuth, index)
}
