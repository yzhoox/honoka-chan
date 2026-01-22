package login

import (
	"encoding/base64"
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	loginschema "honoka-chan/internal/schema/login"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/encrypt"
	honokautils "honoka-chan/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func authKey(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	reqData := gjson.Parse(ctx.MustGet("request_data").(string))
	dummyToken, err := base64.StdEncoding.DecodeString(reqData.Get("dummy_token").String())
	if ss.CheckErr(err) {
		return
	}
	dummyTokenDecrypted := encrypt.RSADecrypt(dummyToken)

	// aesKey := dummyTokenDecrypted[0:16]
	// authData, err := base64.StdEncoding.DecodeString(reqData.Get("auth_data").String())
	// if ss.CheckErr(err) {
	// 	return
	// }
	// authDataDecrypted := encrypt.AESCBCDecrypt(authData, aesKey)[16:]
	// fmt.Println(string(authDataDecrypted))

	clientToken := base64.StdEncoding.EncodeToString(dummyTokenDecrypted)
	serverToken := base64.StdEncoding.EncodeToString([]byte(honokautils.RandomStr(32)))
	authorizeToken := base64.StdEncoding.EncodeToString([]byte(honokautils.RandomStr(32)))

	authJson, err := json.Marshal(map[string]any{
		"client_token": clientToken,
		"server_token": serverToken,
	})
	if ss.CheckErr(err) {
		return
	}

	err = db.Ldb.Set([]byte(authorizeToken), authJson)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(loginschema.AuthKeyResp{
		ResponseData: loginschema.AuthKeyData{
			AuthorizeToken: authorizeToken,
			DummyToken:     serverToken,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/login/authkey", middleware.Common, authKey)
}
