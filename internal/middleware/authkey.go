package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/encrypt"
	honokautils "honoka-chan/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func AuthKey(ctx *gin.Context) {
	req := gjson.Parse(ctx.PostForm("request_data"))
	tDummyToken, err := base64.StdEncoding.DecodeString(req.Get("dummy_token").String())
	utils.CheckErr(err)
	dummyToken := encrypt.RSADecrypt(tDummyToken)

	// aesKey := dummyToken[0:16]
	// tAuthData, err := base64.StdEncoding.DecodeString(req.Get("auth_data").String())
	// utils.CheckErr(err)
	// authData := utils.Sub16(encrypt.AES_CBC_Decrypt(tAuthData, aesKey))
	// fmt.Println(string(authData))

	clientToken := base64.StdEncoding.EncodeToString(dummyToken)
	serverToken := base64.StdEncoding.EncodeToString([]byte(honokautils.RandomStr(32)))
	authorizeToken := base64.StdEncoding.EncodeToString([]byte(honokautils.RandomStr(32)))

	ctx.Set("dummy_token", serverToken)
	ctx.Set("authorize_token", authorizeToken)

	authJson, err := json.Marshal(map[string]any{
		"client_token": clientToken,
		"server_token": serverToken,
	})
	utils.CheckErr(err)
	err = db.DB.Set([]byte(authorizeToken), authJson)
	utils.CheckErr(err)

	nonce := ctx.GetInt("nonce")
	nonce++
	ctx.Set("nonce", nonce)

	authorize := fmt.Sprintf("consumerKey=lovelive_test&timeStamp=%d&version=1.1&nonce=%d&requestTimeStamp=%d", time.Now().Unix(), nonce, ctx.GetInt64("req_time"))

	ctx.Header("user_id", "")
	ctx.Header("authorize", authorize)

	ctx.Next()
}
