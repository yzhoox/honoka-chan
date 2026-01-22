package account

import (
	"encoding/base64"
	"encoding/json"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
)

func loginAuto(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	data, err := ctx.GetRawData()
	if ss.CheckErr(err) {
		return
	}

	data, err = base64.StdEncoding.DecodeString(string(data))
	if ss.CheckErr(err) {
		return
	}

	randKey := ss.GetRandKey()
	decryptedData, err := openssl.Des3ECBDecrypt(data, randKey[0:24], openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	queryStr, _ := url.QueryUnescape(string(decryptedData))
	params, _ := url.ParseQuery(queryStr)
	autoKey := params.Get("autokey")

	var uid, ticket string
	_, err = ss.UserEng.Table("users").Cols("user_id,ticket").Where("autokey = ?", autoKey).Get(&uid, &ticket)
	if ss.CheckErr(err) {
		return
	}

	loginAutoData := ghome.LoginAutoData{}
	loginAutoCode := 0
	loginAutoMsg := "ok"
	if uid != "" {
		loginAutoData = ghome.LoginAutoData{
			Result:  loginAutoCode,
			Message: loginAutoMsg,
			Autokey: autoKey,
			UserId:  uid,
			Ticket:  ticket,
		}
	} else {
		loginAutoCode = 31
		loginAutoMsg = "账号不存在或者登陆状态已过期！"
		loginAutoData = ghome.LoginAutoData{
			Result:  loginAutoCode,
			Message: loginAutoMsg,
		}
	}

	data, err = json.Marshal(loginAutoData)
	if ss.CheckErr(err) {
		return
	}
	encryptedData, err := openssl.Des3ECBEncrypt([]byte(data), randKey[0:24], openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghome.LoginAutoResp{
		Code: loginAutoCode,
		Msg:  loginAutoMsg,
		Data: base64.StdEncoding.EncodeToString(encryptedData),
	})
}

func init() {
	router.AddHandler("v1", "POST", "/account/loginauto", loginAuto)
}
