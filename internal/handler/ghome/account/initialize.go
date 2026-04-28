package account

import (
	"encoding/base64"
	"encoding/json"
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
)

func initialize(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	initData := ghomeschema.InitializeData{
		BrandLogo:                 "http://gskd.sdo.com/ghome/ztc/logo/og/logo_xhdpi.png",
		BrandName:                 "盛趣游戏",
		ForceShowAgreement:        1,
		GreportLogLevel:           "off",
		LogLevel:                  "off",
		LoginButton:               []string{"official"},
		LoginIcon:                 []any{},
		NeedFloatWindowPermission: 1,
		NewDeviceIDServer:         strings.ToUpper(openssl.Md5ToString(ss.GetDeviceID())),
		ShowGuestConfirm:          1,
		VoicetipButton:            1,
	}
	data, err := json.Marshal(initData)
	if ss.CheckErr(err) {
		return
	}
	randKey, err := ss.Get3DESRandKey()
	if ss.CheckErr(err) {
		return
	}
	encryptedData, err := openssl.Des3ECBEncrypt([]byte(data), randKey, openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghomeschema.InitializeResp{
		Code: 0,
		Msg:  "ok",
		Data: base64.StdEncoding.EncodeToString(encryptedData),
	})
}

func init() {
	router.AddHandler("v1", "POST", "/account/initialize", initialize)
}
