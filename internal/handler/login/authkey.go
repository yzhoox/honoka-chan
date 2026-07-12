package login

import (
	"encoding/base64"
	"honoka-chan/internal/middleware"
	loginmodel "honoka-chan/internal/model/login"
	"honoka-chan/internal/router"
	loginschema "honoka-chan/internal/schema/login"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/encrypt"
	utils "honoka-chan/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func authKey(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.FinalizeOrRollback()

	reqData := gjson.Parse(ctx.MustGet("request_data").(string))
	dummyToken, err := base64.StdEncoding.DecodeString(reqData.Get("dummy_token").String())
	if ss.CheckErr(err) {
		return
	}
	dummyTokenDecrypted, err := encrypt.RSADecrypt(dummyToken)
	if ss.CheckErr(err) {
		return
	}

	// aesKey := dummyTokenDecrypted[0:16]
	// authData, err := base64.StdEncoding.DecodeString(reqData.Get("auth_data").String())
	// if ss.CheckErr(err) {
	// 	return
	// }
	// authDataDecrypted := encrypt.AESCBCDecrypt(authData, aesKey)[16:]
	// fmt.Println(string(authDataDecrypted))

	clientToken := base64.StdEncoding.EncodeToString(dummyTokenDecrypted)
	serverToken := base64.StdEncoding.EncodeToString([]byte(utils.RandomStr(32)))
	authorizeToken := base64.StdEncoding.EncodeToString([]byte(utils.RandomStr(32)))

	ss.SetAuthKey(&loginmodel.AuthKey{
		AuthorizeToken: authorizeToken,
		ClientToken:    clientToken,
		ServerToken:    serverToken,
		InsertDate:     time.Now().Format("2006-01-02 15:04:05"),
	})

	ss.Respond(loginschema.AuthKeyResp{
		ResponseData: loginschema.AuthKeyData{
			AuthorizeToken: authorizeToken,
			DummyToken:     serverToken,
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/login/authkey", middleware.Common, authKey)
}
