package middleware

import (
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			recovered := recover()
			if recovered == nil {
				return
			}

			stack := debug.Stack()
			log.Printf("panic serving %s %s: %v\n%s", ctx.Request.Method, ctx.Request.URL.String(), recovered, stack)

			if value, ok := ctx.Get("session"); ok {
				if ss, ok := value.(*session.Session); ok {
					ss.Rollback()
				}
			}

			if !honokautils.IsMainPHPRequest(ctx.Request.URL.Path) {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			honokautils.AbortMaintenanceJSON(ctx, http.StatusInternalServerError, honokautils.NewPanicContent())
		}()

		ctx.Next()
	}
}
