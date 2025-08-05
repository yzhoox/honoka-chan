package middleware

import (
	"fmt"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrorMsg = `{"code":20001,"message":""}`
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Common(ctx *gin.Context) {
	ctx.Set("req_time", time.Now().Unix())

	authorize := ctx.Request.Header.Get("Authorize")
	if authorize == "" {
		ctx.String(http.StatusForbidden, ErrorMsg)
		ctx.Abort()
	}
	ctx.Set("authorize", authorize)

	params, err := url.ParseQuery(authorize)
	utils.CheckErr(err)

	nonce, err := strconv.Atoi(params.Get("nonce"))
	utils.CheckErr(err)
	nonce++
	ctx.Set("nonce", nonce)

	token := params.Get("token")
	ctx.Set("token", token)

	if ctx.Request.URL.String() == "/main.php/login/authkey" ||
		ctx.Request.URL.String() == "/main.php/login/login" {
		// 特殊请求
		fmt.Println("========")
	} else {
		userId := ctx.Request.Header.Get("User-ID")
		if userId == "" {
			ctx.String(http.StatusForbidden, ErrorMsg)
			ctx.Abort()
		}
		ctx.Set("userid", userId)

		rToken, err := db.DB.Get([]byte(userId))
		utils.CheckErr(err)
		if token != string(rToken) {
			ctx.String(http.StatusForbidden, ErrorMsg)
			ctx.Abort()
		}

		if !db.MatchTokenUid(token, userId) {
			ctx.String(http.StatusForbidden, ErrorMsg)
			ctx.Abort()
		}

		ctx.Header("user_id", userId)
		ctx.Header("authorize", fmt.Sprintf("consumerKey=lovelive_test&timeStamp=%d&version=1.1&token=%s&nonce=%d&user_id=%s&requestTimeStamp=%d", time.Now().Unix(), token, nonce, userId, time.Now().Unix()))
	}

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Header("X-Powered-By", "KLab Native APP Platform")
	ctx.Header("server_version", "20120129")
	ctx.Header("Server-Version", "97.4.6")
	ctx.Header("version_up", "0")
	ctx.Header("status_code", "200")

	ctx.Next()
}
