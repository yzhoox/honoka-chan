package middleware

import (
	"fmt"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Common(ctx *gin.Context) {
	reqData := ""
	if form, err := ctx.MultipartForm(); err == nil {
		if v, ok := form.Value["request_data"]; ok && len(v) > 0 {
			reqData = v[0]
		}
	}
	if reqData == "" {
		reqData = ctx.PostForm("request_data")
	}
	ctx.Set("request_data", reqData)

	uid := ctx.GetHeader("User-ID")
	ctx.Set("userid", uid)

	ss := session.Attach(ctx)
	ctx.Set("session", ss)

	if ctx.Request.URL.Path != "/login/authkey" && ctx.Request.URL.Path != "/login/login" {
		if ss.UserPref.ForceRelogin {
			honokautils.AbortMaintenanceJSON(ctx, http.StatusNotFound, honokautils.NewNotFoundContent(ctx.Request.URL.String()))
			return
		}
	}

	authorize := ctx.GetHeader("Authorize")
	params, err := url.ParseQuery(authorize)
	if err != nil {
		honokautils.AbortMaintenanceJSON(ctx, http.StatusNotFound, honokautils.NewNotFoundContent(ctx.Request.URL.String()))
		return
	}

	nonce, _ := strconv.Atoi(params.Get("nonce"))
	nonce++
	ctx.Set("nonce", nonce)

	token := params.Get("token")
	ctx.Set("token", token)

	ctx.Header("user_id", uid)
	ctx.Header("authorize", fmt.Sprintf("consumerKey=lovelive_test&timeStamp=%d&version=1.1&token=%s&nonce=%d&user_id=%s&requestTimeStamp=%d", time.Now().Unix(), token, nonce, uid, time.Now().Unix()))

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Header("X-Powered-By", "KLab Native APP Platform")
	ctx.Header("server_version", "20120129")
	ctx.Header("Server-Version", "97.4.6")
	ctx.Header("version_up", "0")
	ctx.Header("status_code", strconv.Itoa(http.StatusOK))

	ctx.Next()
}
