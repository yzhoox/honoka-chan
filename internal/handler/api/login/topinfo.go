package login

import (
	loginapischema "honoka-chan/internal/schema/api/login"
	"time"
)

func loginTopInfo() (res any, err error) {
	res = loginapischema.TopInfoResp{
		Result: loginapischema.TopInfoData{
			FriendActionCnt:        0,
			FriendGreetCnt:         0,
			FriendVarietyCnt:       0,
			FriendNewCnt:           0,
			PresentCnt:             0,
			SecretBoxBadgeFlag:     false,
			ServerDatetime:         time.Now().Format("2006-01-02 15:04:05"),
			ServerTimestamp:        time.Now().Unix(),
			NoticeFriendDatetime:   time.Now().Format("2006-01-02 15:04:05"),
			NoticeMailDatetime:     "2000-01-01 12:00:00",
			FriendsApprovalWaitCnt: 0,
			FriendsRequestCnt:      0,
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
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
