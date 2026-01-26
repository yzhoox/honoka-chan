package login

import (
	loginapischema "honoka-chan/internal/schema/api/login"
	"time"
)

func loginTopInfoOnce() (res any, err error) {
	res = loginapischema.TopInfoOnceResp{
		Result: loginapischema.TopInfoOnceData{
			NewAchievementCnt:            0,
			UnaccomplishedAchievementCnt: 0,
			LiveDailyRewardExist:         false,
			TrainingEnergy:               10,
			TrainingEnergyMax:            10,
			Notification: loginapischema.Notification{
				Push:       false,
				Lp:         false,
				UpdateInfo: false,
				Campaign:   false,
				Live:       false,
				Lbonus:     false,
				Event:      false,
				Secretbox:  false,
				Birthday:   true,
			},
			OpenArena:               true,
			CostumeStatus:           true,
			OpenAccessory:           true,
			ArenaSiSkillUniqueCheck: true,
			OpenV98:                 true,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
