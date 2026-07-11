package background

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	backgroundschema "honoka-chan/internal/schema/background"
	"honoka-chan/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func set(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	reqData := gjson.Parse(ctx.MustGet("request_data").(string))
	pref := usermodel.UserPref{
		BackgroundID: int(reqData.Get("background_id").Int()),
	}

	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Update(&pref)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(backgroundschema.SetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/background/set", middleware.Common, set)
}
