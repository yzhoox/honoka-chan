package basic

import (
	"encoding/base64"
	"fmt"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/db"
	"honoka-chan/pkg/encrypt"
	honokautils "honoka-chan/pkg/utils"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
)

func handshake(ctx *gin.Context) {
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

	decryptedData := encrypt.RSADecrypt(data)

	params, _ := url.ParseQuery(string(decryptedData))

	randKey := []byte(params.Get("randkey"))
	deviceID := ss.GetDeviceID()

	err = db.Ldb.Set([]byte(deviceID), randKey)
	if ss.CheckErr(err) {
		return
	}

	token := fmt.Sprintf(`{"message":"ok","result":0,"token":"%s"}`, strings.ToUpper(honokautils.RandomStr(33)))
	encryptedToken, err := openssl.Des3ECBEncrypt([]byte(token), randKey[0:24], openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghome.HandshakeResp{
		Code: 0,
		Msg:  "ok",
		Data: base64.StdEncoding.EncodeToString(encryptedToken),
	})
}

func init() {
	router.AddHandler("v1", "POST", "/basic/handshake", handshake)
}
