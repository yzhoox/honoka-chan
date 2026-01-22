package award

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/award"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func set(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	reqData := gjson.Parse(ctx.MustGet("request_data").(string))
	pref := user.UserPref{
		AwardID: int(reqData.Get("award_id").Int()),
	}

	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Update(&pref)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(award.SetResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/award/set", middleware.Common, set)
}
