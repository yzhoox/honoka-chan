package account

import (
	"encoding/base64"
	"encoding/json"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	"net/url"
	"strconv"

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

	randKey, err := ss.Get3DESRandKey()
	if ss.CheckErr(err) {
		return
	}
	decryptedData, err := openssl.Des3ECBDecrypt(data, randKey, openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	queryStr, _ := url.QueryUnescape(string(decryptedData))
	params, _ := url.ParseQuery(queryStr)
	autoKey := params.Get("autokey")

	var userData usermodel.Users
	exists, err := ss.UserEng.Table("users").Where("autokey = ?", autoKey).Get(&userData)
	if ss.CheckErr(err) {
		return
	}

	loginAutoData := ghomeschema.LoginAutoData{}
	loginAutoCode := 0
	loginAutoMsg := "ok"
	if exists {
		loginAutoData = ghomeschema.LoginAutoData{
			Result:  loginAutoCode,
			Message: loginAutoMsg,
			Autokey: autoKey,
			UserId:  strconv.Itoa(userData.UserID),
			Ticket:  userData.Ticket,
		}
	} else {
		loginAutoCode = 31
		loginAutoMsg = "账号不存在或者登陆状态已过期！"
		loginAutoData = ghomeschema.LoginAutoData{
			Result:  loginAutoCode,
			Message: loginAutoMsg,
		}
	}

	data, err = json.Marshal(loginAutoData)
	if ss.CheckErr(err) {
		return
	}
	encryptedData, err := openssl.Des3ECBEncrypt([]byte(data), randKey, openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghomeschema.LoginAutoResp{
		Code: loginAutoCode,
		Msg:  loginAutoMsg,
		Data: base64.StdEncoding.EncodeToString(encryptedData),
	})
}

func init() {
	router.AddHandler("v1", "POST", "/account/loginauto", loginAuto)
}
