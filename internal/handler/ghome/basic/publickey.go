package basic

import (
	"encoding/base64"
	"encoding/pem"
	"honoka-chan/config"
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/pkg/utils"

	"github.com/gin-gonic/gin"
)

func publicKey(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	publicKeyCode := 0
	publicKeyMsg := "ok"
	publicKeyData := ghomeschema.PublicKeyData{
		Result:  publicKeyCode,
		Message: publicKeyMsg,
	}

	publicKey := honokautils.ReadAllText(config.PublicKeyPath)
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		publicKeyMsg = "公钥读取失败！"
		publicKeyCode = 31
	} else {
		publicKeyData.Key = base64.StdEncoding.EncodeToString(block.Bytes)
		publicKeyData.Method = "rsa"
	}

	ss.Respond(ghomeschema.PublicKeyResp{
		Code: publicKeyCode,
		Msg:  publicKeyMsg,
		Data: publicKeyData,
	})
}

func init() {
	router.AddHandler("v1", "POST", "/basic/publickey", publicKey)
}
