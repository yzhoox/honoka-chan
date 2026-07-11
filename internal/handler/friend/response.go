package friend

import (
	"errors"
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	friendschema "honoka-chan/internal/schema/friend"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func response(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	req := friendschema.ResponseReq{}
	err := honokautils.ParseRequestData(ctx, &req)
	if ss.CheckErr(err) {
		return
	}

	targetUserID, err := resolveActualFriendUserID(ss, req.UserID)
	if ss.CheckErr(err) {
		return
	}

	switch req.Status {
	case 0:
		_, err = ss.UserEng.Table(new(usermodel.UserFriend)).
			Where("user_id = ?", ss.UserID).
			Where("friend_user_id = ?", targetUserID).
			Where("status = ?", usermodel.FriendStatusAwaitingApproval).
			Delete()
		if ss.CheckErr(err) {
			return
		}

		_, err = ss.UserEng.Table(new(usermodel.UserFriend)).
			Where("user_id = ?", targetUserID).
			Where("friend_user_id = ?", ss.UserID).
			Where("status = ?", usermodel.FriendStatusPending).
			Delete()
		if ss.CheckErr(err) {
			return
		}
	case 2:
		err = usermodel.EnsureMutualFriendWithIsNew(ss.UserEng, ss.UserID, targetUserID, usermodel.FriendStatusApproved, true, true)
		if ss.CheckErr(err) {
			return
		}
	default:
		ss.CheckErr(errors.New("invalid friend response status"))
		return
	}

	ss.Respond(friendschema.ResponseResp{
		ResponseData: []any{},
		ReleaseInfo:  []any{},
		StatusCode:   http.StatusOK,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/friend/response", middleware.Common, response)
}
