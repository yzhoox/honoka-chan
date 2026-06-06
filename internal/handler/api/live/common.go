package live

import (
	"fmt"
	liveapischema "honoka-chan/internal/schema/api/live"
	"honoka-chan/internal/session"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	randomLiveStartDate = "1970-01-01 00:00:00"
	randomLiveEndDate   = "2038-01-19 12:14:07"
)

var jst = time.FixedZone("JST", 9*3600)

type liveAvailabilityRow struct {
	LiveDifficultyID      int    `xorm:"live_difficulty_id"`
	LiveTrackID           int    `xorm:"live_track_id"`
	LiveSettingID         int    `xorm:"live_setting_id"`
	Difficulty            int    `xorm:"difficulty"`
	AcFlag                int    `xorm:"ac_flag"`
	NotesSettingAsset     string `xorm:"notes_setting_asset"`
	ExcludeClearCountFlag int    `xorm:"exclude_clear_count_flag"`
}

type specialLiveRotationRow struct {
	RotationGroupID  int    `xorm:"rotation_group_id"`
	LiveDifficultyID int    `xorm:"live_difficulty_id"`
	BaseDate         string `xorm:"base_date"`
}

func listAvailableNormalLiveDifficultyIDs(ss *session.Session) ([]int, error) {
	rows, err := listAvailableNormalLiveRows(ss)
	if err != nil {
		return nil, err
	}

	result := make([]int, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.LiveDifficultyID)
	}
	return result, nil
}

func listAvailableNormalLiveTrackIDs(ss *session.Session) ([]int, error) {
	rows, err := listAvailableNormalLiveRows(ss)
	if err != nil {
		return nil, err
	}

	trackSet := map[int]struct{}{}
	for _, row := range rows {
		trackSet[row.LiveTrackID] = struct{}{}
	}

	result := make([]int, 0, len(trackSet))
	for id := range trackSet {
		result = append(result, id)
	}
	sort.Ints(result)
	return result, nil
}

func listAvailableNormalLiveRows(ss *session.Session) ([]liveAvailabilityRow, error) {
	rows := []liveAvailabilityRow{}
	err := ss.MainEng.Table("normal_live_m").Alias("live").
		Join("LEFT", "live_setting_m setting", "live.live_setting_id = setting.live_setting_id").
		Select(`
			live.live_difficulty_id,
			live.live_setting_id,
			setting.live_track_id,
			setting.difficulty,
			setting.ac_flag,
			setting.notes_setting_asset,
			0 AS exclude_clear_count_flag
		`).
		OrderBy("live.live_difficulty_id ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	result := make([]liveAvailabilityRow, 0, len(rows))
	for _, row := range rows {
		if !hasBeatmap(row.NotesSettingAsset) {
			continue
		}
		result = append(result, row)
	}
	return result, nil
}

func listAvailableSpecialLiveRows(ss *session.Session) ([]liveAvailabilityRow, error) {
	rows := []liveAvailabilityRow{}
	err := ss.MainEng.Table("special_live_m").Alias("live").
		Join("LEFT", "live_setting_m setting", "live.live_setting_id = setting.live_setting_id").
		Select(`
			live.live_difficulty_id,
			live.live_setting_id,
			setting.live_track_id,
			setting.difficulty,
			setting.ac_flag,
			setting.notes_setting_asset,
			COALESCE(live.exclude_clear_count_flag, 0) AS exclude_clear_count_flag
		`).
		OrderBy("live.live_difficulty_id ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	result := make([]liveAvailabilityRow, 0, len(rows))
	for _, row := range rows {
		if !hasBeatmap(row.NotesSettingAsset) {
			continue
		}
		result = append(result, row)
	}
	return result, nil
}

func hasBeatmap(notesSettingAsset string) bool {
	if notesSettingAsset == "" {
		return false
	}
	path := filepath.Join("assets/serverdata/beatmaps", notesSettingAsset)
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func listTodaySpecialRotationDifficultyIDs(ss *session.Session, now time.Time) ([]int, error) {
	rows := []specialLiveRotationRow{}
	err := ss.MainEng.Table("special_live_rotation_m").
		OrderBy("rotation_group_id ASC, base_date ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	grouped := map[int][]struct {
		liveDifficultyID int
		dayModulo        int64
	}{}
	for _, row := range rows {
		baseDate, err := parseRotationBaseDate(row.BaseDate)
		if err != nil {
			return nil, err
		}
		grouped[row.RotationGroupID] = append(grouped[row.RotationGroupID], struct {
			liveDifficultyID int
			dayModulo        int64
		}{
			liveDifficultyID: row.LiveDifficultyID,
			dayModulo:        baseDate.Unix() / 86400,
		})
	}

	currentDay := now.In(jst).Unix() / 86400
	groupIDs := make([]int, 0, len(grouped))
	for groupID := range grouped {
		groupIDs = append(groupIDs, groupID)
	}
	sort.Ints(groupIDs)

	result := make([]int, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		liveList := grouped[groupID]
		if len(liveList) == 0 {
			continue
		}
		currentDayModulo := currentDay % int64(len(liveList))
		for _, live := range liveList {
			if live.dayModulo%int64(len(liveList)) == currentDayModulo {
				result = append(result, live.liveDifficultyID)
				break
			}
		}
	}
	return result, nil
}

func listAllSpecialRotationDifficultyIDs(ss *session.Session) (map[int]struct{}, error) {
	ids := []int{}
	err := ss.MainEng.Table("special_live_rotation_m").
		Cols("live_difficulty_id").
		Find(&ids)
	if err != nil {
		return nil, err
	}

	result := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		result[id] = struct{}{}
	}
	return result, nil
}

func listAvailableTrainingLiveDifficultyIDs(ss *session.Session) ([]int, error) {
	normalTrackIDs, err := listAvailableNormalLiveTrackIDs(ss)
	if err != nil {
		return nil, err
	}

	specialRows, err := listAvailableSpecialLiveRows(ss)
	if err != nil {
		return nil, err
	}

	rotationIDs, err := listAllSpecialRotationDifficultyIDs(ss)
	if err != nil {
		return nil, err
	}

	normalTrackIDSet := make(map[int]struct{}, len(normalTrackIDs))
	for _, id := range normalTrackIDs {
		normalTrackIDSet[id] = struct{}{}
	}

	result := make([]int, 0)
	for _, row := range specialRows {
		if _, ok := normalTrackIDSet[row.LiveTrackID]; !ok {
			continue
		}
		if row.Difficulty <= 5 || row.AcFlag != 0 || row.ExcludeClearCountFlag != 0 {
			continue
		}
		if _, ok := rotationIDs[row.LiveDifficultyID]; ok {
			continue
		}
		result = append(result, row.LiveDifficultyID)
	}

	sort.Ints(result)
	return result, nil
}

func parseRotationBaseDate(value string) (time.Time, error) {
	layouts := []string{
		"2006/01/02 15:04:05",
		"2006/01/02 3:04:05",
	}
	for _, layout := range layouts {
		t, err := time.ParseInLocation(layout, value, jst)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid special live rotation base_date: %s", value)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func nextDayStart(t time.Time) time.Time {
	return startOfDay(t).Add(24 * time.Hour)
}

func buildRandomLiveList() []liveapischema.RandomLiveList {
	return []liveapischema.RandomLiveList{
		{AttributeID: 1, StartDate: randomLiveStartDate, EndDate: randomLiveEndDate},
		{AttributeID: 2, StartDate: randomLiveStartDate, EndDate: randomLiveEndDate},
		{AttributeID: 3, StartDate: randomLiveStartDate, EndDate: randomLiveEndDate},
	}
}

func buildNormalLiveStatusList(ids []int, snapshotMap map[int]session.LiveStatusSnapshot) []liveapischema.NormalLiveStatusList {
	result := make([]liveapischema.NormalLiveStatusList, 0, len(ids))
	for _, id := range ids {
		snapshot := snapshotMap[id]
		result = append(result, liveapischema.NormalLiveStatusList{
			LiveDifficultyID:   id,
			Status:             snapshot.Status,
			HiScore:            snapshot.HiScore,
			HiComboCount:       snapshot.HiComboCount,
			ClearCnt:           snapshot.ClearCnt,
			AchievedGoalIDList: snapshot.AchievedGoalIDList,
		})
	}
	return result
}

func buildSpecialLiveStatusList(ids []int, snapshotMap map[int]session.LiveStatusSnapshot) []liveapischema.SpecialLiveStatusList {
	result := make([]liveapischema.SpecialLiveStatusList, 0, len(ids))
	for _, id := range ids {
		snapshot := snapshotMap[id]
		result = append(result, liveapischema.SpecialLiveStatusList{
			LiveDifficultyID:   id,
			Status:             snapshot.Status,
			HiScore:            snapshot.HiScore,
			HiComboCount:       snapshot.HiComboCount,
			ClearCnt:           snapshot.ClearCnt,
			AchievedGoalIDList: snapshot.AchievedGoalIDList,
		})
	}
	return result
}

func buildTrainingLiveStatusList(ids []int, snapshotMap map[int]session.LiveStatusSnapshot) []liveapischema.TrainingLiveStatusList {
	result := make([]liveapischema.TrainingLiveStatusList, 0, len(ids))
	for _, id := range ids {
		snapshot := snapshotMap[id]
		result = append(result, liveapischema.TrainingLiveStatusList{
			LiveDifficultyID:   id,
			Status:             snapshot.Status,
			HiScore:            snapshot.HiScore,
			HiComboCount:       snapshot.HiComboCount,
			ClearCnt:           snapshot.ClearCnt,
			AchievedGoalIDList: snapshot.AchievedGoalIDList,
		})
	}
	return result
}
