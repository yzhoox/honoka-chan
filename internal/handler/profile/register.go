package profile

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	profileschema "honoka-chan/internal/schema/profile"
	"honoka-chan/internal/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func register(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := gjson.Parse(ctx.PostForm("request_data"))
	pref := usermodel.UserPref{
		UserDesc: req.Get("introduction").String(),
	}

	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Update(&pref)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(profileschema.RegisterResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/profile/profileRegister", middleware.Common, register)
}
