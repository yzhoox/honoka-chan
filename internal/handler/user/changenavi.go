package user

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func changeNavi(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := gjson.Parse(ctx.MustGet("request_data").(string))
	pref := user.UserPref{
		UnitOwningUserID: int(req.Get("unit_owning_user_id").Int()),
	}
	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Update(&pref)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(userschema.ChangeNaviResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/user/changeNavi", middleware.Common, changeNavi)
}
