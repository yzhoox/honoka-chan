package login

import (
	"encoding/base64"
	"errors"
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	loginschema "honoka-chan/internal/schema/login"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/encrypt"
	"honoka-chan/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func login(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	has, authKeyData, err := ss.GetAuthKey(ctx.MustGet("token").(string))
	if ss.CheckErr(err) {
		return
	}
	if !has {
		ss.Abort(errors.New("invalid token"))
		return
	}

	clientToken, err := base64.StdEncoding.DecodeString(authKeyData.ClientToken)
	if ss.CheckErr(err) {
		return
	}
	serverToken, err := base64.StdEncoding.DecodeString(authKeyData.ServerToken)
	if ss.CheckErr(err) {
		return
	}

	xmcKey := utils.XOR(clientToken, serverToken)
	aesKey := xmcKey[0:16]

	reqData := gjson.Parse(ctx.MustGet("request_data").(string))
	key, err := base64.StdEncoding.DecodeString(reqData.Get("login_key").String())
	if ss.CheckErr(err) {
		return
	}
	loginKey := encrypt.AESCBCDecrypt(key, aesKey)[16:]
	// fmt.Println("loginKey", string(loginKey))

	// password, err := base64.StdEncoding.DecodeString(reqData.Get("login_passwd").String())
	// if ss.CheckErr(err) {
	// 	return
	// }
	// loginPasswd := encrypt.AESCBCDecrypt(password, aesKey)[16:]
	// fmt.Println("loginPasswd", string(loginPasswd))

	authorizeToken := base64.StdEncoding.EncodeToString([]byte(utils.RandomStr(32)))

	var userID int
	exists, err := ss.UserEng.Table("user_key").Where("key = ?", string(loginKey)).Cols("user_id").Get(&userID)
	if ss.CheckErr(err) {
		return
	}

	if !exists || userID == 0 {
		ss.Abort(errors.New("invalid user"))
		return
	}

	if err := usermodel.ClearUserForceRelogin(ss.UserEng, userID); ss.CheckErr(err) {
		return
	}

	ss.Respond(loginschema.LoginResp{
		ResponseData: loginschema.LoginData{
			AuthorizeToken:  authorizeToken,
			UserId:          userID,
			ServerTimestamp: time.Now().Unix(),
			AdultFlag:       2,
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/login/login", middleware.Common, login)
}
