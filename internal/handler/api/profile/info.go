package profile

import (
	"honoka-chan/internal/model/user"
	"honoka-chan/internal/schema/api/profile"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func profileInfo(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	pref := user.UserPref{}
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if err != nil {
		return nil, err
	}

	commonUnit, err := ss.MainEng.Table("common_unit_m").Count()
	if err != nil {
		return nil, err
	}

	userUnit, err := ss.UserEng.Table("user_unit").Where("user_id = ?", ss.UserID).Count()
	if err != nil {
		return nil, err
	}

	unitData := profile.UnitData{}
	exists, err := ss.MainEng.Table("common_unit_m").Where("unit_owning_user_id = ?", pref.UnitOwningUserID).Get(&unitData)
	if err != nil {
		return nil, err
	}

	isCommon := true
	if !exists {
		_, err = ss.UserEng.Table("user_unit").
			Where("unit_owning_user_id = ? AND user_id = ?", pref.UnitOwningUserID, ss.UserID).Get(&unitData)
		if err != nil {
			return nil, err
		}
		isCommon = false
	}

	var attrId, maxHp, baseSmile, basePure, baseCool int
	var smileMax, pureMax, coolMax int
	if isCommon {
		// 公共卡片仅为100级属性
		_, err = ss.MainEng.Table("unit_m").Where("unit_id = ?", unitData.UnitID).
			Cols("attribute_id,hp_max,smile_max,pure_max,cool_max").Get(&attrId, &maxHp, &baseSmile, &basePure, &baseCool)
		if err != nil {
			return nil, err
		}

		// 偷懒起见不计算饰品、宝石、回忆画廊等属性加成
		smileMax = baseSmile
		pureMax = basePure
		coolMax = baseCool
		// } else {
		// 	// 用户卡片需要根据等级计算属性
		// 	// TODO
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
	accessoryInfo := profile.AccessoryInfo{
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

	res = profile.InfoResp{
		Result: profile.InfoData{
			UserInfo: profile.UserInfo{
				UserID:               pref.UserID,
				Name:                 pref.UserName,
				Level:                pref.UserLevel,
				CostMax:              100,
				UnitMax:              5000,
				EnergyMax:            417,
				FriendMax:            99,
				UnitCnt:              int(commonUnit + userUnit),
				InviteCode:           pref.InviteCode,
				ElapsedTimeFromLogin: "14\u5c0f\u65f6\u524d",
				Introduction:         pref.UserDesc,
			},
			CenterUnitInfo: profile.CenterUnitInfo{
				UnitOwningUserID:           unitData.UnitOwningUserID,
				UnitID:                     unitData.UnitID,
				Exp:                        unitData.Exp,
				NextExp:                    unitData.NextExp,
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
				Attribute:                  attrId,
				Smile:                      baseSmile,
				Cute:                       basePure,
				Cool:                       baseCool,
				IsLoveMax:                  unitData.IsLoveMax,
				IsLevelMax:                 unitData.IsLevelMax,
				IsRankMax:                  unitData.IsRankMax,
				IsSigned:                   unitData.IsSigned,
				IsSkillLevelMax:            unitData.IsSkillLevelMax,
				SettingAwardID:             pref.AwardID,
				RemovableSkillIds:          removeSkillIds,
				AccessoryInfo:              accessoryInfo,
				Costume:                    profile.Costume{},
				TotalSmile:                 smileMax,
				TotalCute:                  pureMax,
				TotalCool:                  coolMax,
				TotalHp:                    maxHp,
			},
			NaviUnitInfo: profile.NaviUnitInfo{
				UnitOwningUserID:            unitData.UnitOwningUserID,
				UnitID:                      unitData.UnitID,
				Exp:                         unitData.Exp,
				NextExp:                     unitData.NextExp,
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
				TotalSmile:                  smileMax,
				TotalCute:                   pureMax,
				TotalCool:                   coolMax,
				TotalHp:                     maxHp,
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
