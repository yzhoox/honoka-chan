package friend

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	friendschema "honoka-chan/internal/schema/friend"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func expel(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := friendschema.UserIDReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	targetUserID, err := resolveActualFriendUserID(ss, req.UserID)
	if ss.CheckErr(err) {
		return
	}

	_, err = ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("friend_user_id = ?", targetUserID).
		Delete()
	if ss.CheckErr(err) {
		return
	}

	_, err = ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", targetUserID).
		Where("friend_user_id = ?", ss.UserID).
		Delete()
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(friendschema.ExpelResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/friend/expel", middleware.Common, expel)
}
