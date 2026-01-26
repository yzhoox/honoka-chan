package live

import (
	"errors"
	"honoka-chan/internal/middleware"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func partyList(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	pref := usermodel.UserPref{}
	_, err := ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if ss.CheckErr(err) {
		return
	}

	// TODO: 好友功能实装前先使用自己的卡组助战
	var unitList []usermodel.UserDeckUnit
	err = ss.UserEng.Table("user_deck_unit").Where("user_id = ? AND position = 5", ss.UserID).Find(&unitList)
	if ss.CheckErr(err) {
		return
	}

	var partyList []liveschema.PartyList
	for _, u := range unitList {
		var unitInfo unitmodel.UnitDataMap
		has, err := ss.GetBasicUnitInfo().
			Where("unit_owning_user_id = ?", u.UnitOwningUserID).Get(&unitInfo)
		if ss.CheckErr(err) {
			return
		}

		if !has {
			ss.Abort(errors.New("卡片不存在！"))
			return
		}

		partyList = append(partyList, liveschema.PartyList{
			UserInfo: ss.GetUserInfo(),
			CenterUnitInfo: liveschema.CenterUnitInfo{
				UnitOwningUserID:           u.UnitOwningUserID,
				UnitID:                     u.UnitID,
				Exp:                        unitInfo.Exp,
				NextExp:                    0,
				Level:                      u.Level,
				LevelLimitID:               u.LevelLimitID,
				MaxLevel:                   unitInfo.MaxLevel,
				Rank:                       unitInfo.Rank,
				MaxRank:                    unitInfo.MaxRank,
				Love:                       u.Love,
				MaxLove:                    u.MaxLove,
				UnitSkillLevel:             u.UnitSkillLevel,
				MaxHp:                      unitInfo.MaxHp,
				FavoriteFlag:               unitInfo.FavoriteFlag,
				DisplayRank:                u.DisplayRank,
				UnitSkillExp:               unitInfo.UnitSkillExp,
				UnitRemovableSkillCapacity: unitInfo.UnitRemovableSkillCapacity,
				Attribute:                  unitInfo.Attribute,
				Smile:                      unitInfo.Smile,
				Cute:                       unitInfo.Cute,
				Cool:                       unitInfo.Cool,
				IsLoveMax:                  u.IsLoveMax,
				IsLevelMax:                 u.IsLevelMax,
				IsRankMax:                  u.IsRankMax,
				IsSigned:                   u.IsSigned,
				IsSkillLevelMax:            unitInfo.IsSkillLevelMax,
				SettingAwardID:             pref.AwardID,
				RemovableSkillIds:          []int{},
			},
			SettingAwardID:       pref.AwardID,
			AvailableSocialPoint: 10,
			FriendStatus:         1,
		})
	}

	ss.Respond(liveschema.PartyListResp{
		ResponseData: liveschema.PartyListData{
			PartyList:         partyList,
			TrainingEnergy:    10,
			TrainingEnergyMax: 10,
			ServerTimestamp:   time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/live/partyList", middleware.Common, partyList)
}
