package secretbox

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	secretboxschema "honoka-chan/internal/schema/secretbox"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func multi(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	multiReq := secretboxschema.MultiReq{}
	err := honokautils.ParseRequestData(ctx, &multiReq)
	if ss.CheckErr(err) {
		return
	}

	data, err := drawData(ctx, multiReq.SecretBoxID, multiReq.ID)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(secretboxschema.MultiResp{
		ResponseData: data,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/secretbox/multi", middleware.Common, multi)
}
