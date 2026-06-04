package login

import (
	usermodel "honoka-chan/internal/model/user"
	loginapischema "honoka-chan/internal/schema/api/login"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func loginTopInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)
	now := time.Now()

	friendsRequestCnt64, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("status = ?", usermodel.FriendStatusAwaitingApproval).
		Where("is_new = ?", true).
		Count()
	if err != nil {
		return nil, err
	}

	friendsApprovalWaitCnt64, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("status = ?", usermodel.FriendStatusPending).
		Count()
	if err != nil {
		return nil, err
	}

	friendsRequestCnt := int(friendsRequestCnt64)
	friendsApprovalWaitCnt := int(friendsApprovalWaitCnt64)

	res = loginapischema.TopInfoResp{
		Result: loginapischema.TopInfoData{
			FriendActionCnt:        0,
			FriendGreetCnt:         0,
			FriendVarietyCnt:       0,
			FriendNewCnt:           friendsRequestCnt + friendsApprovalWaitCnt,
			PresentCnt:             0,
			SecretBoxBadgeFlag:     false,
			ServerDatetime:         now.Format("2006-01-02 15:04:05"),
			ServerTimestamp:        now.Unix(),
			NoticeFriendDatetime:   now.Format("2006-01-02 15:04:05"),
			NoticeMailDatetime:     "2000-01-01 12:00:00",
			FriendsApprovalWaitCnt: friendsApprovalWaitCnt,
			FriendsRequestCnt:      friendsRequestCnt,
			IsTodayBirthday:        false,
			LicenseInfo: loginapischema.LicenseInfo{
				LicenseList:  []any{},
				LicensedInfo: []any{},
				ExpiredInfo:  []any{},
				BadgeFlag:    false,
			},
			UsingBuffInfo:     []any{},
			IsKlabIDTaskFlag:  false,
			KlabIDTaskCanSync: false,
			HasUnreadAnnounce: false,
			ExchangeBadgeCnt:  []int{0, 0, 0},
			AdFlag:            false,
			HasAdReward:       false,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  now.Unix(),
	}

	return res, err
}
