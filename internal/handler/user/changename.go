package user

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func changeName(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := gjson.Parse(ctx.MustGet("request_data").(string))
	var oldName string
	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Cols("user_name").Get(&oldName)
	if ss.CheckErr(err) {
		return
	}

	pref := usermodel.UserPref{
		UserName: req.Get("name").String(),
	}
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Update(&pref)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(userschema.ChangeNameResp{
		ResponseData: userschema.ChangeNameData{
			BeforeName:      oldName,
			AfterName:       req.Get("name").String(),
			ServerTimestamp: time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/user/changeName", middleware.Common, changeName)
}
