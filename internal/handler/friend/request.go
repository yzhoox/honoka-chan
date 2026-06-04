package friend

import (
	"errors"
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	friendschema "honoka-chan/internal/schema/friend"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

func request(ctx *gin.Context) {
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
	if targetUserID == ss.UserID {
		ss.CheckErr(errors.New("cannot add self as friend"))
		return
	}

	isFriend, err := areUsersFriends(ss, ss.UserID, targetUserID)
	if ss.CheckErr(err) {
		return
	}
	if isFriend {
		ss.Respond(friendschema.RequestResp{
			ResponseData: friendschema.RequestData{
				IsFriend: true,
			},
			ReleaseInfo: []any{},
			StatusCode:  200,
		})
		return
	}

	reversePending, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", targetUserID).
		Where("friend_user_id = ?", ss.UserID).
		Where("status = ?", usermodel.FriendStatusPending).
		Exist()
	if ss.CheckErr(err) {
		return
	}

	selfAwaitingApproval, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("friend_user_id = ?", targetUserID).
		Where("status = ?", usermodel.FriendStatusAwaitingApproval).
		Exist()
	if ss.CheckErr(err) {
		return
	}

	if reversePending && selfAwaitingApproval {
		err = usermodel.EnsureMutualFriendWithIsNew(ss.UserEng, ss.UserID, targetUserID, usermodel.FriendStatusApproved, true, true)
		if ss.CheckErr(err) {
			return
		}

		ss.Respond(friendschema.RequestResp{
			ResponseData: friendschema.RequestData{
				IsFriend: true,
			},
			ReleaseInfo: []any{},
			StatusCode:  200,
		})
		return
	}

	err = usermodel.EnsureFriendLink(ss.UserEng, ss.UserID, targetUserID, usermodel.FriendStatusPending, false)
	if ss.CheckErr(err) {
		return
	}
	err = usermodel.EnsureFriendLink(ss.UserEng, targetUserID, ss.UserID, usermodel.FriendStatusAwaitingApproval, true)
	if ss.CheckErr(err) {
		return
	}

	ss.Respond(friendschema.RequestResp{
		ResponseData: friendschema.RequestData{
			IsFriend: false,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/friend/request", middleware.Common, request)
}
