package live

import (
	"encoding/json"
	"fmt"
	unitmodel "honoka-chan/internal/model/unit"
	usermodel "honoka-chan/internal/model/user"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/utils"
	"math"
)

type liveSettingData struct {
	NotesSettingAsset string `xorm:"notes_setting_asset"`
	ARankScore        int    `xorm:"a_rank_score"`
	BRankScore        int    `xorm:"b_rank_score"`
	CRankScore        int    `xorm:"c_rank_score"`
	SRankScore        int    `xorm:"s_rank_score"`
	AcFlag            int    `xorm:"ac_flag"`
	SwingFlag         int    `xorm:"swing_flag"`
}

type leaderSkillData struct {
	AttributeID      int     `xorm:"attribute_id"`
	MainEffectValue  float64 `xorm:"main_effect_value"`
	ExtraEffectValue float64 `xorm:"extra_effect_value"`
	MemberTagID      int     `xorm:"member_tag_id"`
}

type museumBuff struct {
	Smile float64 `xorm:"smile_buff"`
	Pure  float64 `xorm:"pure_buff"`
	Cool  float64 `xorm:"cool_buff"`
}

type accessoryBonus struct {
	Smile float64 `xorm:"smile_max"`
	Pure  float64 `xorm:"pure_max"`
	Cool  float64 `xorm:"cool_max"`
}

type memberTagMatchKey struct {
	UnitTypeID  int
	MemberTagID int
}

type accessoryBonusRow struct {
	AccessoryOwningUserID int     `xorm:"accessory_owning_user_id"`
	AccessoryID           int     `xorm:"accessory_id"`
	Smile                 float64 `xorm:"smile_max"`
	Pure                  float64 `xorm:"pure_max"`
	Cool                  float64 `xorm:"cool_max"`
}

type removableSkillMeta struct {
	UnitRemovableSkillID int     `xorm:"unit_removable_skill_id"`
	EffectRange          int     `xorm:"effect_range"`
	EffectType           int     `xorm:"effect_type"`
	EffectValue          float64 `xorm:"effect_value"`
	FixedValueFlag       int     `xorm:"fixed_value_flag"`
	TargetReferenceType  int     `xorm:"target_reference_type"`
}

func loadLiveSetting(ss *session.Session, difficultyID int) (liveSettingData, error) {
	setting := liveSettingData{}
	sql := `
SELECT notes_setting_asset,
	a_rank_score,
	b_rank_score,
	c_rank_score,
	s_rank_score,
	ac_flag,
	swing_flag
FROM live_setting_m
WHERE live_setting_id IN (
	SELECT live_setting_id
	FROM normal_live_m
	WHERE live_difficulty_id = ?
	UNION ALL
	SELECT live_setting_id
	FROM special_live_m
	WHERE live_difficulty_id = ?
)
	`
	has, err := ss.MainEng.SQL(sql, difficultyID, difficultyID).Get(&setting)
	if err != nil {
		return liveSettingData{}, err
	}
	if !has {
		return liveSettingData{}, fmt.Errorf("live setting not found: %d", difficultyID)
	}
	return setting, nil
}

func loadLiveNotes(notesSettingAsset string) ([]liveschema.NotesList, error) {
	notes := []liveschema.NotesList{}
	noteData := utils.ReadAllText("./assets/serverdata/beatmaps/" + notesSettingAsset)
	err := json.Unmarshal([]byte(noteData), &notes)
	return notes, err
}

func buildRankInfo(setting liveSettingData) []liveschema.RankInfo {
	return []liveschema.RankInfo{
		{Rank: 5, RankMin: 0, RankMax: setting.CRankScore},
		{Rank: 4, RankMin: setting.CRankScore + 1, RankMax: setting.BRankScore},
		{Rank: 3, RankMin: setting.BRankScore + 1, RankMax: setting.ARankScore},
		{Rank: 2, RankMin: setting.ARankScore + 1, RankMax: setting.SRankScore},
		{Rank: 1, RankMin: setting.SRankScore + 1, RankMax: 0},
	}
}

func getMuseumBuff(ss *session.Session) (museumBuff, error) {
	buff := museumBuff{}
	_, err := ss.MainEng.Table("museum_contents_m").
		Select("SUM(smile_buff) AS smile_buff,SUM(pure_buff) AS pure_buff,SUM(cool_buff) AS cool_buff").
		Get(&buff)
	return buff, err
}

func getDeckCenterUnitID(ss *session.Session, deckID int) (int, error) {
	var centerUnitID int
	_, err := ss.UserEng.Table("user_deck_unit").
		Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
		Where("user_deck.deck_id = ? AND user_deck.user_id = ? AND user_deck_unit.position = 5", deckID, ss.UserID).
		Cols("user_deck_unit.unit_id").
		Get(&centerUnitID)
	return centerUnitID, err
}

func getLeaderSkillData(ss *session.Session, unitID int) (leaderSkillData, error) {
	data := leaderSkillData{}
	_, err := ss.MainEng.Table("unit_m").
		Join("LEFT", "unit_leader_skill_m", "unit_m.default_leader_skill_id = unit_leader_skill_m.unit_leader_skill_id").
		Join("LEFT", "unit_leader_skill_extra_m", "unit_m.default_leader_skill_id = unit_leader_skill_extra_m.unit_leader_skill_id").
		Where("unit_m.unit_id = ?", unitID).
		Cols(`
			unit_m.attribute_id,
			unit_leader_skill_m.effect_value AS main_effect_value,
			unit_leader_skill_extra_m.effect_value AS extra_effect_value,
			unit_leader_skill_extra_m.member_tag_id
		`).
		Get(&data)
	return data, err
}

func applyAttributeBonus(attributeID int, effectValue, smile, pure, cool float64) (float64, float64, float64) {
	switch attributeID {
	case 1:
		return math.Ceil(smile * (effectValue / 100)), 0, 0
	case 2:
		return 0, math.Ceil(pure * (effectValue / 100)), 0
	case 3:
		return 0, 0, math.Ceil(cool * (effectValue / 100))
	default:
		return 0, 0, 0
	}
}

func matchMemberTag(ss *session.Session, cache map[memberTagMatchKey]bool, unitTypeID, memberTagID int) (bool, error) {
	if memberTagID <= 0 {
		return false, nil
	}

	key := memberTagMatchKey{UnitTypeID: unitTypeID, MemberTagID: memberTagID}
	if matched, ok := cache[key]; ok {
		return matched, nil
	}

	matched, err := ss.MainEng.Table("unit_type_member_tag_m").
		Where("unit_type_id = ? AND member_tag_id = ?", unitTypeID, memberTagID).
		Exist()
	if err != nil {
		return false, err
	}
	cache[key] = matched
	return matched, nil
}

func loadDeckUnitDataMap(ss *session.Session, unitOwningUserIDs []int) (map[int]unitmodel.UnitDataMap, error) {
	unitRows := []unitmodel.UnitDataMap{}
	err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", ss.UserID).
		In("a.unit_owning_user_id", unitOwningUserIDs).
		Find(&unitRows)
	if err != nil {
		return nil, err
	}

	unitMap := make(map[int]unitmodel.UnitDataMap, len(unitRows))
	for _, row := range unitRows {
		unitMap[row.UnitOwningUserID] = row
	}
	return unitMap, nil
}

func loadAccessoryBonusMap(ss *session.Session, userID int, unitOwningUserIDs []int) (map[int]accessoryBonus, error) {
	wears := []usermodel.UserAccessoryWear{}
	err := ss.UserEng.Table(new(usermodel.UserAccessoryWear)).
		Where("user_id = ?", userID).
		In("unit_owning_user_id", unitOwningUserIDs).
		Find(&wears)
	if err != nil {
		return nil, err
	}
	if len(wears) == 0 {
		return map[int]accessoryBonus{}, nil
	}

	accessoryIDs := make([]int, 0, len(wears))
	unitToAccessoryID := make(map[int]int, len(wears))
	for _, wear := range wears {
		if wear.AccessoryOwningUserID <= 0 {
			continue
		}
		accessoryIDs = append(accessoryIDs, wear.AccessoryOwningUserID)
		unitToAccessoryID[wear.UnitOwningUserID] = wear.AccessoryOwningUserID
	}
	if len(accessoryIDs) == 0 {
		return map[int]accessoryBonus{}, nil
	}

	rows := []accessoryBonusRow{}
	err = ss.UserEng.Table("user_accessory").
		Where("user_id = ?", userID).
		In("accessory_owning_user_id", accessoryIDs).
		Cols("accessory_owning_user_id,accessory_id").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	accessoryIDSet := make(map[int]struct{}, len(rows))
	for _, row := range rows {
		accessoryIDSet[row.AccessoryID] = struct{}{}
	}

	accessoryIDList := make([]int, 0, len(accessoryIDSet))
	for accessoryID := range accessoryIDSet {
		accessoryIDList = append(accessoryIDList, accessoryID)
	}

	staticRows := []accessoryBonusRow{}
	err = ss.MainEng.Table("accessory_m").
		In("accessory_id", accessoryIDList).
		Cols("accessory_id,smile_max,pure_max,cool_max").
		Find(&staticRows)
	if err != nil {
		return nil, err
	}

	staticBonusMap := make(map[int]accessoryBonus, len(staticRows))
	for _, row := range staticRows {
		staticBonusMap[row.AccessoryID] = accessoryBonus{
			Smile: row.Smile,
			Pure:  row.Pure,
			Cool:  row.Cool,
		}
	}

	accessoryMap := make(map[int]accessoryBonus, len(rows))
	for _, row := range rows {
		bonus, ok := staticBonusMap[row.AccessoryID]
		if !ok {
			continue
		}
		accessoryMap[row.AccessoryOwningUserID] = bonus
	}

	result := make(map[int]accessoryBonus, len(unitToAccessoryID))
	for unitOwningUserID, accessoryOwningUserID := range unitToAccessoryID {
		if bonus, ok := accessoryMap[accessoryOwningUserID]; ok {
			result[unitOwningUserID] = bonus
		}
	}
	return result, nil
}

func loadDeckRemovableSkillMap(ss *session.Session, userID int, unitOwningUserIDs []int) (map[int][]int, error) {
	rows := []usermodel.UserUnitSkillEquip{}
	err := ss.UserEng.Table(new(usermodel.UserUnitSkillEquip)).
		Where("user_id = ?", userID).
		In("unit_owning_user_id", unitOwningUserIDs).
		Find(&rows)
	if err != nil {
		return nil, err
	}

	skillMap := make(map[int][]int, len(unitOwningUserIDs))
	for _, row := range rows {
		skillMap[row.UnitOwningUserID] = append(skillMap[row.UnitOwningUserID], row.UnitRemovableSkillID)
	}
	return skillMap, nil
}

func loadRemovableSkillMetaMap(ss *session.Session, skillIDs []int) (map[int]removableSkillMeta, error) {
	if len(skillIDs) == 0 {
		return map[int]removableSkillMeta{}, nil
	}

	rows := []removableSkillMeta{}
	err := ss.MainEng.Table("unit_removable_skill_m").
		In("unit_removable_skill_id", skillIDs).
		Cols("unit_removable_skill_id,effect_range,effect_type,effect_value,fixed_value_flag,target_reference_type").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	metaMap := make(map[int]removableSkillMeta, len(rows))
	for _, row := range rows {
		metaMap[row.UnitRemovableSkillID] = row
	}
	return metaMap, nil
}
