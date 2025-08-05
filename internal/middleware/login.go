package middleware

import (
	"encoding/base64"
	"honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/encrypt"
	honokautils "honoka-chan/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func Login(ctx *gin.Context) {
	authData, err := db.DB.Get([]byte(ctx.GetString("token")))
	utils.CheckErr(err)

	clientToken, err := base64.StdEncoding.DecodeString(gjson.Get(string(authData), "client_token").String())
	utils.CheckErr(err)
	serverToken, err := base64.StdEncoding.DecodeString(gjson.Get(string(authData), "server_token").String())
	utils.CheckErr(err)

	xmcKey := honokautils.SliceXor(clientToken, serverToken)
	aesKey := xmcKey[0:16]

	req := gjson.Parse(ctx.GetString("request_data"))
	tKey, err := base64.StdEncoding.DecodeString(req.Get("login_key").String())
	utils.CheckErr(err)
	loginKey := honokautils.Sub16(encrypt.AESCBCDecrypt(tKey, aesKey))
	ctx.Set("login_key", string(loginKey))

	tPasswd, err := base64.StdEncoding.DecodeString(req.Get("login_passwd").String())
	utils.CheckErr(err)
	loginPasswd := honokautils.Sub16(encrypt.AESCBCDecrypt(tPasswd, aesKey))
	ctx.Set("login_passwd", string(loginPasswd))

	nonce := ctx.GetInt("nonce")
	nonce++
	ctx.Set("nonce", nonce)

	authorizeToken := base64.StdEncoding.EncodeToString([]byte(honokautils.RandomStr(32)))
	ctx.Set("authorize_token", authorizeToken)

	ctx.Next()
}
