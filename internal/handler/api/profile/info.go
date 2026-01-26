package profile

import (
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	profileapischema "honoka-chan/internal/schema/api/profile"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func profileInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	pref := usermodel.UserPref{}
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if err != nil {
		return nil, err
	}

	unitCount, err := ss.UserEng.Table(new(usermodel.UserUnitData)).
		Where("user_id = ?", ss.UserID).Count()
	if err != nil {
		return nil, err
	}

	unitData := unitmodel.UnitDataMap{}
	_, err = ss.GetBasicUnitInfo().
		Where("a.unit_owning_user_id = ?", pref.UnitOwningUserID).Get(&unitData)
	if err != nil {
		return nil, err
	}

	var accessoryOwningId, accessoryId, exp int
	_, err = ss.UserEng.Table("user_accessory_wear").Where("unit_owning_user_id = ? AND user_id = ?", pref.UnitOwningUserID, ss.UserID).
		Cols("accessory_owning_user_id").Get(&accessoryOwningId)
	if err != nil {
		return nil, err
	}
	_, err = ss.MainEng.Table("common_accessory_m").Where("accessory_owning_user_id = ?", accessoryOwningId).
		Cols("accessory_id,exp").Get(&accessoryId, &exp)
	if err != nil {
		return nil, err
	}
	accessoryInfo := profileapischema.AccessoryInfo{
		AccessoryOwningUserID: accessoryOwningId,
		AccessoryID:           accessoryId,
		Exp:                   exp,
		NextExp:               0,
		Level:                 8,
		MaxLevel:              8,
		RankUpCount:           4,
		FavoriteFlag:          true,
	}

	removeSkillIds := []int{}
	err = ss.UserEng.Table("user_unit_skill_equip").Where("unit_owning_user_id = ? AND user_id = ?", pref.UnitOwningUserID, ss.UserID).
		Cols("unit_removable_skill_id").Find(&removeSkillIds)
	if err != nil {
		return nil, err
	}

	res = profileapischema.InfoResp{
		Result: profileapischema.InfoData{
			UserInfo: profileapischema.UserInfo{
				UserID:               pref.UserID,
				Name:                 pref.UserName,
				Level:                pref.UserLevel,
				CostMax:              100,
				UnitMax:              5000,
				EnergyMax:            417,
				FriendMax:            99,
				UnitCnt:              int(unitCount),
				InviteCode:           pref.InviteCode,
				ElapsedTimeFromLogin: "14\u5c0f\u65f6\u524d",
				Introduction:         pref.UserDesc,
			},
			// TODO: 区分队伍中心卡片和看板卡片
			CenterUnitInfo: profileapischema.CenterUnitInfo{
				UnitOwningUserID:           unitData.UnitOwningUserID,
				UnitID:                     unitData.UnitID,
				Exp:                        unitData.Exp,
				NextExp:                    0,
				Level:                      unitData.Level,
				LevelLimitID:               unitData.LevelLimitID,
				MaxLevel:                   unitData.MaxLevel,
				Rank:                       unitData.Rank,
				MaxRank:                    unitData.MaxRank,
				Love:                       unitData.Love,
				MaxLove:                    unitData.MaxLove,
				UnitSkillLevel:             unitData.UnitSkillLevel,
				MaxHp:                      unitData.MaxHp,
				FavoriteFlag:               unitData.FavoriteFlag,
				DisplayRank:                unitData.DisplayRank,
				UnitSkillExp:               unitData.UnitSkillExp,
				UnitRemovableSkillCapacity: unitData.UnitRemovableSkillCapacity,
				Attribute:                  unitData.Attribute,
				Smile:                      unitData.Smile,
				Cute:                       unitData.Cute,
				Cool:                       unitData.Cool,
				IsLoveMax:                  unitData.IsLoveMax,
				IsLevelMax:                 unitData.IsLevelMax,
				IsRankMax:                  unitData.IsRankMax,
				IsSigned:                   unitData.IsSigned,
				IsSkillLevelMax:            unitData.IsSkillLevelMax,
				SettingAwardID:             pref.AwardID,
				RemovableSkillIds:          removeSkillIds,
				AccessoryInfo:              accessoryInfo,
				Costume:                    profileapischema.Costume{},
				TotalSmile:                 unitData.Smile, // TODO: 加成计算
				TotalCute:                  unitData.Cute,  // 同上
				TotalCool:                  unitData.Cool,  // 同上
				TotalHp:                    unitData.MaxHp,
			},
			NaviUnitInfo: profileapischema.NaviUnitInfo{
				UnitOwningUserID:            unitData.UnitOwningUserID,
				UnitID:                      unitData.UnitID,
				Exp:                         unitData.Exp,
				NextExp:                     0,
				Level:                       unitData.Level,
				MaxLevel:                    unitData.MaxLevel,
				LevelLimitID:                unitData.LevelLimitID,
				Rank:                        unitData.Rank,
				MaxRank:                     unitData.MaxRank,
				Love:                        unitData.Love,
				MaxLove:                     unitData.MaxLove,
				UnitSkillExp:                unitData.UnitSkillExp,
				UnitSkillLevel:              unitData.UnitSkillLevel,
				MaxHp:                       unitData.MaxHp,
				UnitRemovableSkillCapacity:  unitData.UnitRemovableSkillCapacity,
				FavoriteFlag:                unitData.FavoriteFlag,
				DisplayRank:                 unitData.DisplayRank,
				IsLoveMax:                   unitData.IsLoveMax,
				IsLevelMax:                  unitData.IsLevelMax,
				IsRankMax:                   unitData.IsRankMax,
				IsSigned:                    unitData.IsSigned,
				IsSkillLevelMax:             unitData.IsSkillLevelMax,
				IsRemovableSkillCapacityMax: unitData.IsRemovableSkillCapacityMax,
				InsertDate:                  "2016-10-11 10:33:03",
				TotalSmile:                  unitData.Smile, // TODO: 加成计算
				TotalCute:                   unitData.Cute,  // 同上
				TotalCool:                   unitData.Cool,  // 同上
				TotalHp:                     unitData.MaxHp,
				RemovableSkillIds:           removeSkillIds,
			},
			IsAlliance:          false,
			FriendStatus:        0,
			SettingAwardID:      pref.AwardID,
			SettingBackgroundID: pref.BackgroundID,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
