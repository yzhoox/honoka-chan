package session

import (
	usermodel "honoka-chan/internal/model/user"
	liveschema "honoka-chan/internal/schema/live"
	liverecordschema "honoka-chan/internal/schema/liverecord"
)

// Move from LevelDB to SQLite
// https://github.com/YumeMichi/honoka-chan/blob/a749efad9fd0789668dedcb59948248da369d32c/handler/live.go#L554
// https://github.com/DarkEnergyProcessor/NPPS4/blob/29aaba6e7a3b4b80d414e58f4b511a9627cb0e24/npps4/system/live.py#L360
func (ss *Session) GetLiveInProgress() (bool, *usermodel.UserLiveInProgress) {
	progress := usermodel.UserLiveInProgress{}
	has, err := ss.UserEng.Table(new(usermodel.UserLiveInProgress)).
		Where("user_id = ?", ss.UserID).Get(&progress)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &progress
}

// Move from LevelDB to SQLite
// https://github.com/YumeMichi/honoka-chan/blob/a749efad9fd0789668dedcb59948248da369d32c/handler/live.go#L45
// https://github.com/DarkEnergyProcessor/NPPS4/blob/29aaba6e7a3b4b80d414e58f4b511a9627cb0e24/npps4/system/live.py#L372
func (ss *Session) RegisterLiveInProgress(deckID int) {
	var err error
	has, progress := ss.GetLiveInProgress()
	if has {
		progress.DeckID = deckID
		_, err = ss.UserEng.Table(new(usermodel.UserLiveInProgress)).ID(progress.ID).Update(progress)
	} else {
		progress = &usermodel.UserLiveInProgress{
			DeckID: deckID,
			UserID: ss.UserID,
		}
		_, err = ss.UserEng.Table(new(usermodel.UserLiveInProgress)).Insert(progress)
	}
	if ss.CheckErr(err) {
		return
	}
}

// Move from LevelDB to SQLite
// https://github.com/YumeMichi/honoka-chan/blob/a749efad9fd0789668dedcb59948248da369d32c/handler/live.go#L554
// https://github.com/DarkEnergyProcessor/NPPS4/blob/29aaba6e7a3b4b80d414e58f4b511a9627cb0e24/npps4/system/live.py#L366
func (ss *Session) ClearLiveInProgress() {
	_, err := ss.UserEng.Table(new(usermodel.UserLiveInProgress)).Where("user_id = ?", ss.UserID).Delete()
	if ss.CheckErr(err) {
		return
	}
}

func (ss *Session) GetLiveInfo(LiveDifficultyID int) (bool, *liveschema.LiveInfo) {
	var info struct {
		LiveDifficultyID  int    `xorm:"live_difficulty_id"`
		NotesSettingAsset string `xorm:"notes_setting_asset"`
		ARankScore        int    `xorm:"a_rank_score"`
		BRankScore        int    `xorm:"b_rank_score"`
		CRankScore        int    `xorm:"c_rank_score"`
		SRankScore        int    `xorm:"s_rank_score"`
		ARankCombo        int    `xorm:"a_rank_combo"`
		BRankCombo        int    `xorm:"b_rank_combo"`
		CRankCombo        int    `xorm:"c_rank_combo"`
		SRankCombo        int    `xorm:"s_rank_combo"`
		ARankComplete     int    `xorm:"a_rank_complete"`
		BRankComplete     int    `xorm:"b_rank_complete"`
		CRankComplete     int    `xorm:"c_rank_complete"`
		SRankComplete     int    `xorm:"s_rank_complete"`
		AcFlag            int    `xorm:"ac_flag"`
		SwingFlag         int    `xorm:"swing_flag"`
	}
	has, err := ss.MainEng.Table("special_live_m").Alias("a").
		Join("LEFT", "live_setting_m", "a.live_setting_id = live_setting_m.live_setting_id").
		Where("live_difficulty_id = ?", LiveDifficultyID).
		Select(`
			a.live_difficulty_id,
			a.c_rank_complete,
			a.b_rank_complete,
			a.a_rank_complete,
			a.s_rank_complete,
			live_setting_m.notes_setting_asset,
			live_setting_m.a_rank_score,
			live_setting_m.b_rank_score,
			live_setting_m.c_rank_score,
			live_setting_m.s_rank_score,
			live_setting_m.a_rank_combo,
			live_setting_m.b_rank_combo,
			live_setting_m.c_rank_combo,
			live_setting_m.s_rank_combo,
			live_setting_m.ac_flag,
			live_setting_m.swing_flag
		`).Get(&info)
	if ss.CheckErr(err) {
		return false, nil
	}

	if !has {
		has, err = ss.MainEng.Table("normal_live_m").Alias("a").
			Join("LEFT", "live_setting_m", "a.live_setting_id = live_setting_m.live_setting_id").
			Where("live_difficulty_id = ?", LiveDifficultyID).
			Select(`
				a.live_difficulty_id,
				a.c_rank_complete,
				a.b_rank_complete,
				a.a_rank_complete,
				a.s_rank_complete,
				live_setting_m.notes_setting_asset,
				live_setting_m.a_rank_score,
				live_setting_m.b_rank_score,
				live_setting_m.c_rank_score,
				live_setting_m.s_rank_score,
				live_setting_m.a_rank_combo,
				live_setting_m.b_rank_combo,
				live_setting_m.c_rank_combo,
				live_setting_m.s_rank_combo,
				live_setting_m.ac_flag,
				live_setting_m.swing_flag
			`).Get(&info)
		if ss.CheckErr(err) {
			return false, nil
		}
	}

	// var notesList []liveschema.NotesList
	// if has {
	// 	noteData := utils.ReadAllText("./assets/serverdata/beatmaps/" + info.NotesSettingAsset)
	// 	err = json.Unmarshal([]byte(noteData), &notesList)
	// 	if ss.CheckErr(err) {
	// 		return false, nil
	// 	}
	// }

	return has, &liveschema.LiveInfo{
		LiveDifficultyID: info.LiveDifficultyID,
		IsRandom:         false,
		ARankScore:       info.ARankScore,
		BRankScore:       info.BRankScore,
		CRankScore:       info.CRankScore,
		SRankScore:       info.SRankScore,
		ARankCombo:       info.ARankCombo,
		BRankCombo:       info.BRankCombo,
		CRankCombo:       info.CRankCombo,
		SRankCombo:       info.SRankCombo,
		ARankComplete:    info.ARankComplete,
		BRankComplete:    info.BRankComplete,
		CRankComplete:    info.CRankComplete,
		SRankComplete:    info.SRankComplete,
		AcFlag:           info.AcFlag,
		SwingFlag:        info.SwingFlag,
		NotesList:        []liveschema.NotesList{},
	}
}

func (ss *Session) GetDeckInfo(deckID int) (bool, *liverecordschema.DeckInfo) {
	// 卡片信息
	totalHp, totalSmile, totalCute, totalCool := 0, 0, 0, 0
	unitData := []liverecordschema.UnitList{}
	units := ss.GetUserDeckUnit(deckID)
	for _, u := range units {
		_, uData := ss.GetUnitInfo(u.UnitID)
		tempUnitData := liverecordschema.UnitList{
			UnitID:                     u.UnitID,
			Position:                   u.Position,
			Level:                      u.Level,
			LevelLimitID:               u.LevelLimitID,
			DisplayRank:                u.DisplayRank,
			Love:                       u.MaxLove,
			UnitSkillLevel:             u.UnitSkillLevel,
			IsRankMax:                  u.IsRankMax,
			IsLoveMax:                  u.IsLoveMax,
			IsLevelMax:                 u.IsLevelMax,
			IsSigned:                   u.IsSigned,
			Rank:                       uData.Rank,
			MaxHp:                      uData.MaxHp,
			UnitRemovableSkillCapacity: uData.UnitRemovableSkillCapacity,
			TotalStatus: liverecordschema.TotalStatus{
				Hp:    uData.MaxHp,
				Smile: uData.Smile,
				Cute:  uData.Cute,
				Cool:  uData.Cool,
			},
			SiBonus: liverecordschema.SiBonus{
				Hp:    0,
				Smile: 0,
				Cute:  0,
				Cool:  0,
			},
		}

		// 技能宝石信息
		skillData := ss.GetUserUnitSkillEquipID(u.UnitOwningUserID)
		tempUnitData.RemovableSkillIds = skillData

		// 饰品信息
		has, accessoryData := ss.GetUserAccessoryInfoByUnitOwningUserID(u.UnitOwningUserID)
		if has {
			tempUnitData.AccessoryInfo = accessoryData
		}

		// 其他属性
		totalHp += uData.MaxHp
		totalSmile += uData.Smile
		totalCute += uData.Cute
		totalCool += uData.Cool

		unitData = append(unitData, tempUnitData)
	}

	return true, &liverecordschema.DeckInfo{
		TotalStatus: liverecordschema.TotalStatus{
			Hp:    totalHp,
			Smile: totalSmile,
			Cute:  totalCute,
			Cool:  totalCool,
		},
		CenterBonus: liverecordschema.CenterBonus{ // TODO: 计算C位加成
			Hp:    0,
			Smile: 0,
			Cute:  0,
			Cool:  0,
		},
		SiBonus: liverecordschema.SiBonus{
			Hp:    0,
			Smile: 0,
			Cute:  0,
			Cool:  0,
		},
		UnitList: unitData,
	}
}
