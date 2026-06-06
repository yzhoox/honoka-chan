package login

import (
	loginapischema "honoka-chan/internal/schema/api/login"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func loginTopInfoOnce(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)
	userInfo := ss.GetUserInfo()

	res = loginapischema.TopInfoOnceResp{
		Result: loginapischema.TopInfoOnceData{
			NewAchievementCnt:            0,
			UnaccomplishedAchievementCnt: 0,
			LiveDailyRewardExist:         false,
			TrainingEnergy:               userInfo.TrainingEnergy,
			TrainingEnergyMax:            userInfo.TrainingEnergyMax,
			Notification: loginapischema.Notification{
				Push:       false,
				Lp:         false,
				UpdateInfo: false,
				Campaign:   false,
				Live:       false,
				Lbonus:     false,
				Event:      false,
				Secretbox:  false,
				Birthday:   ss.UserPref.HasBirthDate(),
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
