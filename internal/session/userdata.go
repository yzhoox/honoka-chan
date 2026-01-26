package session

import (
	"errors"
	"honoka-chan/internal/constant"
	usermodel "honoka-chan/internal/model/user"
	userschema "honoka-chan/internal/schema/user"
)

func (ss *Session) GetUserPref(userID string) usermodel.UserPref {
	pref := usermodel.UserPref{}
	has, err := ss.UserEng.Table(new(usermodel.UserPref)).
		Where("user_id = ?", userID).Get(&pref)
	if ss.CheckErr(err) {
		return usermodel.UserPref{}
	}
	if !has {
		ss.Abort(errors.New("用户不存在！"))
		return usermodel.UserPref{}
	}

	return pref
}

func (ss *Session) GetUserInfo() userschema.UserInfo {
	return userschema.UserInfo{
		UserID:                         ss.UserID,
		Name:                           ss.UserPref.UserName,
		Level:                          ss.UserPref.UserLevel,
		Exp:                            1089696,
		PreviousExp:                    0,
		NextExp:                        1207185,
		GameCoin:                       112124104,
		SnsCoin:                        10000,
		FreeSnsCoin:                    50000,
		PaidSnsCoin:                    50000,
		SocialPoint:                    1438165,
		UnitMax:                        5000,
		WaitingUnitMax:                 5000,
		CurrentEnergy:                  417,
		EnergyMax:                      417,
		EnergyFullTime:                 "2023-03-20 01:28:55",
		LicenseLiveEnergyRecoverlyTime: 60,
		EnergyFullNeedTime:             0,
		OverMaxEnergy:                  417,
		TrainingEnergy:                 10,
		TrainingEnergyMax:              10,
		FriendMax:                      99,
		InviteCode:                     ss.UserPref.InviteCode,
		InsertDate:                     "2015-08-10 18:58:30",
		UpdateDate:                     "2018-08-09 18:13:12",
		TutorialState:                  -1,
		DiamondCoin:                    0,
		CrystalCoin:                    0,
		UnlockRandomLiveMuse:           1,
		UnlockRandomLiveAqours:         1,
	}
}

func (ss *Session) GetUserLiveRecord(liveDifficultyID int) []usermodel.UserLiveRecord {
	liveRecord := []usermodel.UserLiveRecord{}
	err := ss.UserEng.Table(new(usermodel.UserLiveRecord)).
		Where("user_id = ? AND live_difficulty_id = ?",
			ss.UserID, liveDifficultyID).
		Find(&liveRecord)
	if ss.CheckErr(err) {
		return []usermodel.UserLiveRecord{}
	}

	return liveRecord
}

func (ss *Session) GetUserLiveRecordWithSkillOn(liveDifficultyID int, isSkillOn bool) (bool, *usermodel.UserLiveRecord) {
	liveRecord := usermodel.UserLiveRecord{}
	has, err := ss.UserEng.Table(new(usermodel.UserLiveRecord)).
		Where("user_id = ? AND live_difficulty_id = ? AND is_skill_on = ?",
			ss.UserID, liveDifficultyID, isSkillOn).
		Get(&liveRecord)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &liveRecord
}

func (ss *Session) UpdateUserLiveRecord(newLiveRecord *usermodel.UserLiveRecord, updateType constant.PreciseScoreUpdateType) {
	// 获取原有记录
	has, liveRecord := ss.GetUserLiveRecordWithSkillOn(newLiveRecord.LiveDifficultyID, newLiveRecord.IsSkillOn)
	if !has {
		// 没有记录则直接新增
		ss.UpdateUserLiveRecordWrapper(newLiveRecord, false)
		return
	}

	// 判断记录更新类型
	switch updateType {
	case constant.PreciseScoreUpdateTypePerfect:
		// 根据 Perfect 数更新
		// 判断 Perfect 数是否大于原有记录
		if newLiveRecord.PerfectCnt >= liveRecord.PerfectCnt {
			// 更新记录
			ss.UpdateUserLiveRecordWrapper(newLiveRecord, true)
		}
	case constant.PreciseScoreUpdateTypeScore:
		// 根据分数更新
		// 判断分数是否大于原有记录
		if newLiveRecord.TotalScore >= liveRecord.TotalScore {
			// 更新记录
			ss.UpdateUserLiveRecordWrapper(newLiveRecord, true)
		}
	default:
		// 总是更新
		ss.UpdateUserLiveRecordWrapper(newLiveRecord, true)
	}
}

func (ss *Session) UpdateUserLiveRecordWrapper(newLiveRecord *usermodel.UserLiveRecord, isUpdate bool) {
	if isUpdate {
		_, err := ss.UserEng.Table(new(usermodel.UserLiveRecord)).
			Where("user_id = ? AND live_difficulty_id = ? AND is_skill_on = ?",
				ss.UserID, newLiveRecord.LiveDifficultyID, newLiveRecord.IsSkillOn).
			Update(newLiveRecord)
		if ss.CheckErr(err) {
			return
		}
	} else {
		_, err := ss.UserEng.Insert(newLiveRecord)
		if ss.CheckErr(err) {
			return
		}
	}
}
