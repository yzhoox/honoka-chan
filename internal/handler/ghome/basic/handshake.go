package basic

import (
	"encoding/base64"
	"errors"
	"fmt"
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/encrypt"
	"honoka-chan/pkg/utils"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-think/openssl"
)

func handshake(ctx *gin.Context) {
	ss := session.Attach(ctx)
	if ss.Done() {
		return
	}
	defer ss.FinalizeOrRollback()

	data, err := ctx.GetRawData()
	if ss.CheckErr(err) {
		return
	}

	data, err = base64.StdEncoding.DecodeString(string(data))
	if ss.CheckErr(err) {
		return
	}

	decryptedData := encrypt.RSADecrypt(data)

	params, err := url.ParseQuery(string(decryptedData))
	if ss.CheckErr(err) {
		return
	}

	randKey := params.Get("randkey")
	if len(randKey) < 24 {
		ss.Abort(errors.New("invalid rand key"))
		return
	}
	ss.SetRandKey(randKey)

	token := fmt.Sprintf(`{"message":"ok","result":0,"token":"%s"}`, strings.ToUpper(utils.RandomStr(33)))
	encryptedToken, err := openssl.Des3ECBEncrypt([]byte(token), []byte(randKey[:24]), openssl.PKCS7_PADDING)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghomeschema.HandshakeResp{
		Code: 0,
		Msg:  "ok",
		Data: base64.StdEncoding.EncodeToString(encryptedToken),
	})
}

func init() {
	router.AddHandler("v1", "POST", "/basic/handshake", handshake)
}
