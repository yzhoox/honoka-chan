package account

import (
	"encoding/base64"
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
)

func reportRole(ctx *gin.Context) {
	ss := session.Attach(ctx)
	defer ss.FinalizeOrRollback()

	randKey, err := ss.Get3DESRandKey()
	if ss.CheckErr(err) {
		return
	}
	token := `{"message":"ok"}`
	encryptedToken, err := openssl.Des3ECBEncrypt([]byte(token), randKey, openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghomeschema.ReportRoleResp{
		Code: 0,
		Msg:  "ok",
		Data: base64.StdEncoding.EncodeToString(encryptedToken),
	})
}

func init() {
	router.AddHandler("v1", "POST", "/account/reportRole", reportRole)
}
