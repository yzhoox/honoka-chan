package basic

import (
	"encoding/json"
	"honoka-chan/internal/router"
	ghomeschema "honoka-chan/internal/schema/ghome"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func getCode(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	codeArray := `{"codeArray":[{"btntext":"好的","code":"-10264022","msg_from":2,"text":"","title":"短信验证码被阻止","type":1},{"btntext":"","code":"-10869623","msg_from":2,"text":"","title":"网络连接失败，无法一键登录","type":2},{"btntext":"","code":"10298300","msg_from":2,"text":"","title":"","type":3},{"btntext":"","code":"10298311","msg_from":2,"text":"","title":"","type":3},{"btntext":"","code":"10298312","msg_from":2,"text":"","title":"","type":3},{"btntext":"","code":"10298313","msg_from":2,"text":"","title":"","type":1},{"btntext":"","code":"10298321","msg_from":2,"text":"","title":"","type":3},{"btntext":"","code":"10298322","msg_from":2,"text":"","title":"","type":3}],"codeVersion":"1.0.5"}`
	getCodeData := map[string]any{}
	err := json.Unmarshal([]byte(codeArray), &getCodeData)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(ghomeschema.GetCodeResp{
		Code: 0,
		Msg:  "ok",
		Data: getCodeData,
	})
}

func init() {
	router.AddHandler("v1", "GET", "/basic/getcode", getCode)
	router.AddHandler("v1", "POST", "/basic/getcode", getCode)
}
