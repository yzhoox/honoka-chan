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

	friendsRequestCnt, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("status = ?", usermodel.FriendStatusAwaitingApproval).
		Where("is_new = ?", true).
		Count()
	if err != nil {
		return nil, err
	}

	friendsApprovalWaitCnt, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("status = ?", usermodel.FriendStatusPending).
		Count()
	if err != nil {
		return nil, err
	}

	friendGreetCnt, err := ss.UserEng.Table(new(usermodel.UserGreet)).
		Where("receiver_id = ?", ss.UserID).
		Where("deleted_from_receiver = ?", false).
		Where("readed = ?", false).
		Count()
	if err != nil {
		return nil, err
	}

	res = loginapischema.TopInfoResp{
		Result: loginapischema.TopInfoData{
			FriendActionCnt:        friendGreetCnt,
			FriendGreetCnt:         friendGreetCnt,
			FriendVarietyCnt:       0,
			FriendNewCnt:           friendsRequestCnt + friendsApprovalWaitCnt,
			PresentCnt:             0,
			SecretBoxBadgeFlag:     false,
			ServerDatetime:         now.Format("2006-01-02 15:04:05"),
			ServerTimestamp:        now.Unix(),
			NoticeFriendDatetime:   now.Format("2006-01-02 15:04:05"),
			NoticeMailDatetime:     now.Format("2006-01-02 15:04:05"),
			FriendsApprovalWaitCnt: friendsApprovalWaitCnt,
			FriendsRequestCnt:      friendsRequestCnt,
			IsTodayBirthday:        ss.UserPref.IsBirthdayToday(now),
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
