package live

import (
	"fmt"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	liveschema "honoka-chan/internal/schema/live"
	userschema "honoka-chan/internal/schema/user"
	"honoka-chan/internal/session"
)

type liveSupportRow struct {
	FriendUserID           int    `xorm:"friend_user_id"`
	UserName               string `xorm:"user_name"`
	UserLevel              int    `xorm:"user_level"`
	InviteCode             string `xorm:"invite_code"`
	AwardID                int    `xorm:"award_id"`
	CenterUnitOwningUserID int    `xorm:"center_unit_owning_user_id"`
}

func listLiveSupportRows(ss *session.Session) ([]liveSupportRow, error) {
	rows := []liveSupportRow{}
	err := ss.UserEng.Table(new(usermodel.UserFriend)).Alias("uf").
		Join("LEFT", "user_pref up", "up.user_id = uf.friend_user_id").
		Join("LEFT", "user_deck ud", "ud.user_id = uf.friend_user_id AND ud.main_flag = 1").
		Join("LEFT", "user_deck_unit udu", "udu.user_deck_id = ud.id AND udu.position = 5").
		Where("uf.user_id = ?", ss.UserID).
		Where("uf.status = ?", usermodel.FriendStatusApproved).
		Select(`
			uf.friend_user_id,
			up.user_name,
			up.user_level,
			up.invite_code,
			up.award_id,
			COALESCE(udu.unit_owning_user_id, up.unit_owning_user_id) AS center_unit_owning_user_id
		`).
		OrderBy("uf.update_date DESC, uf.id DESC").
		Find(&rows)
	return rows, err
}

func buildLiveSupportParty(ss *session.Session, row liveSupportRow) (liveschema.PartyList, error) {
	centerUnitInfo, err := buildLiveCenterUnitInfo(ss, row.FriendUserID, row.CenterUnitOwningUserID, row.AwardID)
	if err != nil {
		return liveschema.PartyList{}, err
	}

	userInfo := buildLiveSupportUserInfo(ss, row)
	return liveschema.PartyList{
		UserInfo:             userInfo,
		CenterUnitInfo:       centerUnitInfo,
		SettingAwardID:       row.AwardID,
		AvailableSocialPoint: 10,
		FriendStatus:         usermodel.ClientFriendStatusApproved,
	}, nil
}

func buildLiveSupportUserInfo(ss *session.Session, row liveSupportRow) userschema.UserInfo {
	userInfo := ss.GetUserInfo()
	userInfo.UserID = row.FriendUserID
	userInfo.Name = row.UserName
	userInfo.Level = row.UserLevel
	userInfo.InviteCode = row.InviteCode
	return userInfo
}

func buildLiveCenterUnitInfo(ss *session.Session, targetUserID, unitOwningUserID, awardID int) (liveschema.CenterUnitInfo, error) {
	info := liveschema.CenterUnitInfo{
		SettingAwardID:    awardID,
		RemovableSkillIds: []int{},
	}
	if unitOwningUserID <= 0 {
		return info, nil
	}

	unitData := unitmodel.UnitDataMap{}
	has, err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", targetUserID).
		Where("a.unit_owning_user_id = ?", unitOwningUserID).
		Get(&unitData)
	if err != nil {
		return info, err
	}
	if !has {
		return info, nil
	}

	removableSkillIDs := []int{}
	err = ss.UserEng.Table(new(usermodel.UserUnitSkillEquip)).
		Where("user_id = ? AND unit_owning_user_id = ?", targetUserID, unitOwningUserID).
		Cols("unit_removable_skill_id").
		Find(&removableSkillIDs)
	if err != nil {
		return info, err
	}

	return liveschema.CenterUnitInfo{
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
		RemovableSkillIds:          removableSkillIDs,
	}, nil
}

func getSupportCenterUnitID(ss *session.Session, partyUserID int, fallbackUnitID int) (int, error) {
	if partyUserID <= 0 || partyUserID == ss.UserID {
		return fallbackUnitID, nil
	}

	exists, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", ss.UserID).
		Where("friend_user_id = ?", partyUserID).
		Where("status = ?", usermodel.FriendStatusApproved).
		Exist()
	if err != nil {
		return 0, err
	}
	if !exists {
		return fallbackUnitID, nil
	}

	var centerUnitID int
	has, err := ss.UserEng.Table(new(usermodel.UserDeckUnit)).Alias("udu").
		Join("LEFT", "user_deck ud", "ud.id = udu.user_deck_id").
		Where("udu.user_id = ?", partyUserID).
		Where("ud.main_flag = ?", 1).
		Where("udu.position = ?", 5).
		Cols("udu.unit_id").
		Get(&centerUnitID)
	if err != nil {
		return 0, err
	}
	if !has || centerUnitID <= 0 {
		return fallbackUnitID, nil
	}

	return centerUnitID, nil
}

func mustFindCenterUnitInfo(ss *session.Session, userID, unitOwningUserID int) (liveschema.CenterUnitInfo, error) {
	info, err := buildLiveCenterUnitInfo(ss, userID, unitOwningUserID, 0)
	if err != nil {
		return liveschema.CenterUnitInfo{}, err
	}
	if info.UnitOwningUserID <= 0 || info.UnitID <= 0 {
		return liveschema.CenterUnitInfo{}, fmt.Errorf("support center unit not found for user %d", userID)
	}
	return info, nil
}
