package user

import (
	"honoka-chan/internal/model/user"
	userapischema "honoka-chan/internal/schema/api/user"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func userInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	pref := user.UserPref{}
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if err != nil {
		return nil, err
	}

	res = userapischema.InfoResp{
		Result: userapischema.UserInfo{
			UserID:                         ss.UserID,
			Name:                           pref.UserName,
			Level:                          pref.UserLevel,
			Exp:                            1089696,
			PreviousExp:                    0,
			NextExp:                        1207185,
			GameCoin:                       999999999,
			SnsCoin:                        100000,
			FreeSnsCoin:                    50000,
			PaidSnsCoin:                    50000,
			SocialPoint:                    1438395,
			UnitMax:                        5000,
			WaitingUnitMax:                 1000,
			EnergyMax:                      417,
			EnergyFullTime:                 "2023-03-20 03:58:55",
			LicenseLiveEnergyRecoverlyTime: 60,
			EnergyFullNeedTime:             0,
			OverMaxEnergy:                  417,
			TrainingEnergy:                 100,
			TrainingEnergyMax:              100,
			FriendMax:                      99,
			InviteCode:                     pref.InviteCode,
			InsertDate:                     "2015-08-10 18:58:30",
			UpdateDate:                     "2018-08-09 18:13:12",
			TutorialState:                  -1,
			DiamondCoin:                    0,
			CrystalCoin:                    0,
			LpRecoveryItem:                 []userapischema.LpRecoveryItem{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
