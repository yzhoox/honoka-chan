package notice

import (
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	noticeschema "honoka-chan/internal/schema/notice"
	"honoka-chan/internal/session"
	"strconv"
	"time"
)

type greetingPeerRow struct {
	UserID                 int    `xorm:"user_id"`
	UserName               string `xorm:"user_name"`
	UserLevel              int    `xorm:"user_level"`
	AwardID                int    `xorm:"award_id"`
	CenterUnitOwningUserID int    `xorm:"center_unit_owning_user_id"`
}

func getGreetingPeer(ss *session.Session, userID int) (noticeschema.GreetingPeer, error) {
	row := greetingPeerRow{}
	has, err := ss.UserEng.Table(new(usermodel.UserPref)).Alias("up").
		Join("LEFT", "user_deck ud", "ud.user_id = up.user_id AND ud.main_flag = 1").
		Join("LEFT", "user_deck_unit udu", "udu.user_deck_id = ud.id AND udu.position = 5").
		Select(`
			up.user_id,
			up.user_name,
			up.user_level,
			up.award_id,
			COALESCE(udu.unit_owning_user_id, up.unit_owning_user_id) AS center_unit_owning_user_id
		`).
		Where("up.user_id = ?", userID).
		Get(&row)
	if err != nil {
		return noticeschema.GreetingPeer{}, err
	}
	if !has {
		return noticeschema.GreetingPeer{}, nil
	}

	centerUnitInfo, err := buildGreetingCenterUnitInfo(ss, userID, row.CenterUnitOwningUserID, row.AwardID)
	if err != nil {
		return noticeschema.GreetingPeer{}, err
	}

	return noticeschema.GreetingPeer{
		UserData: noticeschema.GreetingUserData{
			UserID: row.UserID,
			Name:   row.UserName,
			Level:  row.UserLevel,
		},
		CenterUnitInfo: centerUnitInfo,
		SettingAwardID: row.AwardID,
	}, nil
}

func buildGreetingCenterUnitInfo(ss *session.Session, userID, unitOwningUserID, awardID int) (noticeschema.GreetingCenterUnitInfo, error) {
	info := noticeschema.GreetingCenterUnitInfo{
		SettingAwardID:    awardID,
		RemovableSkillIds: []int{},
	}
	if unitOwningUserID <= 0 {
		return info, nil
	}

	unitData := unitmodel.UnitDataMap{}
	has, err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", userID).
		Where("a.unit_owning_user_id = ?", unitOwningUserID).
		Get(&unitData)
	if err != nil {
		return info, err
	}
	if !has {
		return info, nil
	}

	accessoryInfo, err := getGreetingAccessoryInfo(ss, userID, unitOwningUserID)
	if err != nil {
		return info, err
	}

	removableSkillIDs := []int{}
	err = ss.UserEng.Table(new(usermodel.UserUnitSkillEquip)).
		Where("user_id = ? AND unit_owning_user_id = ?", userID, unitOwningUserID).
		Cols("unit_removable_skill_id").
		Find(&removableSkillIDs)
	if err != nil {
		return info, err
	}

	return noticeschema.GreetingCenterUnitInfo{
		UnitOwningUserID:           int64(unitData.UnitOwningUserID),
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
		RemovableSkillIds:          removableSkillIDs,
		AccessoryInfo:              accessoryInfo,
	}, nil
}

func getGreetingAccessoryInfo(ss *session.Session, userID, unitOwningUserID int) (*noticeschema.GreetingAccessoryInfo, error) {
	accessoryWear := usermodel.UserAccessoryWear{}
	has, err := ss.UserEng.Table(new(usermodel.UserAccessoryWear)).
		Where("user_id = ? AND unit_owning_user_id = ?", userID, unitOwningUserID).
		Get(&accessoryWear)
	if err != nil {
		return nil, err
	}
	if !has || accessoryWear.AccessoryOwningUserID <= 0 {
		return nil, nil
	}

	has, accessoryData := ss.GetAccessoryByAccessoryOwningUserIDForUser(userID, accessoryWear.AccessoryOwningUserID)
	if !has {
		return nil, nil
	}

	return &noticeschema.GreetingAccessoryInfo{
		AccessoryOwningUserID: accessoryWear.AccessoryOwningUserID,
		AccessoryID:           accessoryData.AccessoryID,
		Exp:                   accessoryData.Exp,
		NextExp:               0,
		Level:                 accessoryData.MaxLevel,
		MaxLevel:              accessoryData.MaxLevel,
		RankUpCount:           accessoryData.MaxLevel - accessoryData.DefaultMaxLevel,
		FavoriteFlag:          true,
	}, nil
}

func formatGreetingElapsedTime(ts int64) string {
	if ts <= 0 {
		return "刚刚"
	}

	delta := time.Since(time.Unix(ts, 0))
	if delta < time.Minute {
		return "刚刚"
	}
	if delta < time.Hour {
		return strconv.Itoa(int(delta.Minutes())) + "分钟前"
	}
	if delta < 24*time.Hour {
		return strconv.Itoa(int(delta.Hours())) + "小时前"
	}
	return strconv.Itoa(int(delta.Hours())/24) + "天前"
}
