package profile

import (
	"errors"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	profileapischema "honoka-chan/internal/schema/api/profile"
	"honoka-chan/internal/session"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func profileInfo(ctx *gin.Context, targetUserID int) (res any, err error) {
	ss := session.Get(ctx)

	targetPref, err := getTargetUserPref(ss, targetUserID)
	if err != nil {
		return nil, err
	}
	targetUserID = targetPref.UserID

	unitCount, err := ss.UserEng.Table(new(usermodel.UserUnitData)).
		Where("user_id = ?", targetUserID).
		Count()
	if err != nil {
		return nil, err
	}

	lastLoginTime, err := getTargetLastLoginTime(ss, targetUserID)
	if err != nil {
		return nil, err
	}

	centerUnitOwningUserID, err := getMainDeckCenterUnitOwningUserID(ss, targetUserID, targetPref.UnitOwningUserID)
	if err != nil {
		return nil, err
	}

	centerUnitInfo, centerRaw, err := getProfileCenterUnitInfo(ss, targetUserID, centerUnitOwningUserID, targetPref.AwardID)
	if err != nil {
		return nil, err
	}

	naviUnitInfo, naviRaw, err := getProfileNaviUnitInfo(ss, targetUserID, targetPref.UnitOwningUserID)
	if err != nil {
		return nil, err
	}

	friendStatus, err := usermodel.ResolveClientFriendStatus(ss.UserEng, ss.UserID, targetUserID)
	if err != nil {
		return nil, err
	}

	if naviRaw != nil {
		centerUnitInfo.Costume = profileapischema.Costume{
			UnitID:    naviRaw.UnitID,
			IsRankMax: naviRaw.IsRankMax,
			IsSigned:  naviRaw.IsSigned,
		}
	} else if centerRaw != nil {
		centerUnitInfo.Costume = profileapischema.Costume{
			UnitID:    centerRaw.UnitID,
			IsRankMax: centerRaw.IsRankMax,
			IsSigned:  centerRaw.IsSigned,
		}
	}

	res = profileapischema.InfoResp{
		Result: profileapischema.InfoData{
			UserInfo: profileapischema.UserInfo{
				UserID:               targetPref.UserID,
				Name:                 targetPref.UserName,
				Level:                targetPref.UserLevel,
				CostMax:              100,
				UnitMax:              5000,
				EnergyMax:            targetPref.EffectiveEnergyMax(),
				FriendMax:            99,
				UnitCnt:              unitCount,
				InviteCode:           targetPref.InviteCode,
				ElapsedTimeFromLogin: formatProfileElapsedTime(lastLoginTime),
				Introduction:         targetPref.UserDesc,
			},
			CenterUnitInfo:      centerUnitInfo,
			NaviUnitInfo:        naviUnitInfo,
			IsAlliance:          false,
			FriendStatus:        friendStatus,
			SettingAwardID:      targetPref.AwardID,
			SettingBackgroundID: targetPref.BackgroundID,
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, nil
}

func getTargetUserPref(ss *session.Session, targetUserID int) (usermodel.UserPref, error) {
	if targetUserID <= 0 || targetUserID == ss.UserID {
		return ss.UserPref, nil
	}

	pref := usermodel.UserPref{}
	has, err := ss.UserEng.Table(new(usermodel.UserPref)).
		Where("user_id = ?", targetUserID).
		Get(&pref)
	if err != nil {
		return usermodel.UserPref{}, err
	}
	if !has {
		has, err = ss.UserEng.Table(new(usermodel.UserPref)).
			Where("invite_code = ?", strconv.Itoa(targetUserID)).
			Get(&pref)
		if err != nil {
			return usermodel.UserPref{}, err
		}
	}
	if !has {
		return usermodel.UserPref{}, errors.New("user not found")
	}
	if pref.NeedsProfileMigration() {
		pref.ApplyProfileDefaults()
		_, err = ss.UserEng.Table(new(usermodel.UserPref)).
			Where("user_id = ?", pref.UserID).
			Cols(usermodel.UserPrefProfileColumns()...).
			Update(&pref)
		if err != nil {
			return usermodel.UserPref{}, err
		}
	}

	return pref, nil
}

func getTargetLastLoginTime(ss *session.Session, targetUserID int) (int64, error) {
	if targetUserID == ss.UserID {
		user := usermodel.Users{}
		has, err := ss.UserEng.Table(new(usermodel.Users)).
			Where("user_id = ?", targetUserID).
			Get(&user)
		if err != nil {
			return 0, err
		}
		if !has {
			return 0, errors.New("user not found")
		}
		return user.LastLoginTime, nil
	}

	user := usermodel.Users{}
	has, err := ss.UserEng.Table(new(usermodel.Users)).
		Where("user_id = ?", targetUserID).
		Get(&user)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errors.New("user not found")
	}

	return user.LastLoginTime, nil
}

func getMainDeckCenterUnitOwningUserID(ss *session.Session, targetUserID, fallback int) (int, error) {
	var unitOwningUserID int
	has, err := ss.UserEng.Table(new(usermodel.UserDeckUnit)).Alias("udu").
		Join("LEFT", "user_deck ud", "ud.id = udu.user_deck_id").
		Where("udu.user_id = ?", targetUserID).
		Where("ud.main_flag = ?", 1).
		Where("udu.position = ?", 5).
		Cols("udu.unit_owning_user_id").
		Get(&unitOwningUserID)
	if err != nil {
		return 0, err
	}
	if has && unitOwningUserID > 0 {
		return unitOwningUserID, nil
	}
	return fallback, nil
}

func getProfileCenterUnitInfo(ss *session.Session, targetUserID, unitOwningUserID, awardID int) (profileapischema.CenterUnitInfo, *unitmodel.UnitDataMap, error) {
	unitData, err := getTargetUnitData(ss, targetUserID, unitOwningUserID)
	if err != nil {
		return profileapischema.CenterUnitInfo{}, nil, err
	}
	if unitData == nil {
		return profileapischema.CenterUnitInfo{}, nil, nil
	}

	accessoryInfo, err := getAccessoryInfo(ss, targetUserID, unitOwningUserID)
	if err != nil {
		return profileapischema.CenterUnitInfo{}, nil, err
	}

	removeSkillIDs, err := getRemovableSkillIDs(ss, targetUserID, unitOwningUserID)
	if err != nil {
		return profileapischema.CenterUnitInfo{}, nil, err
	}

	info := profileapischema.CenterUnitInfo{
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
		SettingAwardID:             awardID,
		RemovableSkillIds:          removeSkillIDs,
		AccessoryInfo:              accessoryInfo,
		TotalSmile:                 unitData.Smile,
		TotalCute:                  unitData.Cute,
		TotalCool:                  unitData.Cool,
		TotalHp:                    unitData.MaxHp,
	}

	return info, unitData, nil
}

func getProfileNaviUnitInfo(ss *session.Session, targetUserID, unitOwningUserID int) (profileapischema.NaviUnitInfo, *unitmodel.UnitDataMap, error) {
	unitData, err := getTargetUnitData(ss, targetUserID, unitOwningUserID)
	if err != nil {
		return profileapischema.NaviUnitInfo{}, nil, err
	}
	if unitData == nil {
		return profileapischema.NaviUnitInfo{}, nil, nil
	}

	removeSkillIDs, err := getRemovableSkillIDs(ss, targetUserID, unitOwningUserID)
	if err != nil {
		return profileapischema.NaviUnitInfo{}, nil, err
	}

	info := profileapischema.NaviUnitInfo{
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
		IsRankMax:                   unitData.IsRankMax,
		IsLoveMax:                   unitData.IsLoveMax,
		IsLevelMax:                  unitData.IsLevelMax,
		IsSigned:                    unitData.IsSigned,
		IsSkillLevelMax:             unitData.IsSkillLevelMax,
		IsRemovableSkillCapacityMax: unitData.IsRemovableSkillCapacityMax,
		InsertDate:                  "2016-10-11 10:33:03",
		TotalSmile:                  unitData.Smile,
		TotalCute:                   unitData.Cute,
		TotalCool:                   unitData.Cool,
		TotalHp:                     unitData.MaxHp,
		RemovableSkillIds:           removeSkillIDs,
	}

	return info, unitData, nil
}

func getTargetUnitData(ss *session.Session, targetUserID, unitOwningUserID int) (*unitmodel.UnitDataMap, error) {
	if unitOwningUserID <= 0 {
		return nil, nil
	}

	unitData := unitmodel.UnitDataMap{}
	has, err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", targetUserID).
		Where("a.unit_owning_user_id = ?", unitOwningUserID).
		Get(&unitData)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}

	return &unitData, nil
}

func getAccessoryInfo(ss *session.Session, targetUserID, unitOwningUserID int) (*profileapischema.AccessoryInfo, error) {
	if unitOwningUserID <= 0 {
		return nil, nil
	}

	var accessoryOwningID int
	has, err := ss.UserEng.Table("user_accessory_wear").
		Where("unit_owning_user_id = ? AND user_id = ?", unitOwningUserID, targetUserID).
		Cols("accessory_owning_user_id").
		Get(&accessoryOwningID)
	if err != nil {
		return nil, err
	}
	if !has || accessoryOwningID <= 0 {
		return nil, nil
	}

	has, accessoryData := ss.GetAccessoryByAccessoryOwningUserID(accessoryOwningID)
	if !has {
		return nil, nil
	}

	return &profileapischema.AccessoryInfo{
		AccessoryOwningUserID: accessoryOwningID,
		AccessoryID:           accessoryData.AccessoryID,
		Exp:                   accessoryData.Exp,
		NextExp:               0,
		Level:                 accessoryData.MaxLevel,
		MaxLevel:              accessoryData.MaxLevel,
		RankUpCount:           accessoryData.MaxLevel - accessoryData.DefaultMaxLevel,
		FavoriteFlag:          true,
	}, nil
}

func getRemovableSkillIDs(ss *session.Session, targetUserID, unitOwningUserID int) ([]int, error) {
	removeSkillIDs := []int{}
	if unitOwningUserID <= 0 {
		return removeSkillIDs, nil
	}

	err := ss.UserEng.Table("user_unit_skill_equip").
		Where("unit_owning_user_id = ? AND user_id = ?", unitOwningUserID, targetUserID).
		Cols("unit_removable_skill_id").
		Find(&removeSkillIDs)
	if err != nil {
		return nil, err
	}

	return removeSkillIDs, nil
}

func formatProfileElapsedTime(ts int64) string {
	if ts <= 0 {
		return "刚刚"
	}

	delta := time.Since(time.Unix(ts, 0))
	if delta < time.Minute {
		return "刚刚"
	}
	if delta < time.Hour {
		return timeDurationInt(delta.Minutes()) + "分钟前"
	}
	if delta < 24*time.Hour {
		return timeDurationInt(delta.Hours()) + "小时前"
	}
	return timeDurationInt(delta.Hours()/24) + "天前"
}

func timeDurationInt(value float64) string {
	return strconv.Itoa(int(value))
}
