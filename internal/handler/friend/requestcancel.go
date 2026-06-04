package friend

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	friendschema "honoka-chan/internal/schema/friend"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func requestCancel(ctx *gin.Context) {
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
		Where("status = ?", usermodel.FriendStatusPending).
		Delete()
	if ss.CheckErr(err) {
		return
	}

	_, err = ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", targetUserID).
		Where("friend_user_id = ?", ss.UserID).
		Where("status = ?", usermodel.FriendStatusAwaitingApproval).
		Delete()
	if ss.CheckErr(err) {
		return
	}

	isFriend, err := areUsersFriends(ss, ss.UserID, targetUserID)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(friendschema.RequestCancelResp{
		ResponseData: friendschema.RequestCancelData{
			IsFriend: isFriend,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/friend/requestCancel", middleware.Common, requestCancel)
}
