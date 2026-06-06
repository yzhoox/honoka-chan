package live

import (
	"fmt"
	"honoka-chan/config"
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

var cst = time.FixedZone("CST", 8*3600)

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

type specialLiveRotationSchedule struct {
	LiveDifficultyID int
	StartTime        time.Time
	EndTime          time.Time
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
			dayModulo:        dayIndexInCST(baseDate),
		})
	}

	currentDay := dayIndexInCST(now)
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

func listCurrentAndNextSpecialRotationSchedules(ss *session.Session, now time.Time) ([]specialLiveRotationSchedule, error) {
	if config.Conf.Settings.UnlockAllSpecialRotation {
		return listAllSpecialRotationSchedules(ss)
	}

	rows := []specialLiveRotationRow{}
	err := ss.MainEng.Table("special_live_rotation_m").
		OrderBy("rotation_group_id ASC, base_date ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	type rotationEntry struct {
		liveDifficultyID int
		baseTime         time.Time
		dayModulo        int64
	}

	grouped := map[int][]rotationEntry{}
	groupIDs := []int{}
	for _, row := range rows {
		baseDate, err := parseRotationBaseDate(row.BaseDate)
		if err != nil {
			return nil, err
		}
		if _, ok := grouped[row.RotationGroupID]; !ok {
			groupIDs = append(groupIDs, row.RotationGroupID)
		}
		grouped[row.RotationGroupID] = append(grouped[row.RotationGroupID], rotationEntry{
			liveDifficultyID: row.LiveDifficultyID,
			baseTime:         baseDate,
			dayModulo:        dayIndexInCST(baseDate),
		})
	}
	sort.Ints(groupIDs)

	currentDay := dayIndexInCST(now)
	result := make([]specialLiveRotationSchedule, 0, len(groupIDs)*2)
	for _, groupID := range groupIDs {
		liveList := grouped[groupID]
		if len(liveList) == 0 {
			continue
		}

		currentIndex := -1
		currentDayModulo := currentDay % int64(len(liveList))
		for i, live := range liveList {
			if live.dayModulo%int64(len(liveList)) == currentDayModulo {
				currentIndex = i
				break
			}
		}
		if currentIndex == -1 {
			continue
		}

		periodDays := int64(1)
		if len(liveList) >= 2 {
			diff := int64(liveList[1].baseTime.Sub(liveList[0].baseTime) / (24 * time.Hour))
			if diff > 0 {
				periodDays = diff
			}
		}

		currentEntry := liveList[currentIndex]
		offsetDays := (currentDay - currentEntry.dayModulo) % periodDays
		if offsetDays < 0 {
			offsetDays += periodDays
		}
		currentStartTime := startOfDay(now.In(cst)).AddDate(0, 0, -int(offsetDays))
		nextStartTime := currentStartTime.AddDate(0, 0, int(periodDays))
		duration := time.Duration(periodDays)*24*time.Hour - time.Second
		nextIndex := (currentIndex + 1) % len(liveList)

		for i, live := range liveList {
			switch i {
			case currentIndex:
				result = append(result, specialLiveRotationSchedule{
					LiveDifficultyID: live.liveDifficultyID,
					StartTime:        currentStartTime,
					EndTime:          currentStartTime.Add(duration),
				})
			case nextIndex:
				result = append(result, specialLiveRotationSchedule{
					LiveDifficultyID: live.liveDifficultyID,
					StartTime:        nextStartTime,
					EndTime:          nextStartTime.Add(duration),
				})
			}
		}
	}

	return result, nil
}

func listAllSpecialRotationSchedules(ss *session.Session) ([]specialLiveRotationSchedule, error) {
	rows := []specialLiveRotationRow{}
	err := ss.MainEng.Table("special_live_rotation_m").
		OrderBy("rotation_group_id ASC, base_date ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	type rotationEntry struct {
		liveDifficultyID int
		baseTime         time.Time
	}

	grouped := map[int][]rotationEntry{}
	groupIDs := []int{}
	for _, row := range rows {
		baseDate, err := parseRotationBaseDate(row.BaseDate)
		if err != nil {
			return nil, err
		}
		if _, ok := grouped[row.RotationGroupID]; !ok {
			groupIDs = append(groupIDs, row.RotationGroupID)
		}
		grouped[row.RotationGroupID] = append(grouped[row.RotationGroupID], rotationEntry{
			liveDifficultyID: row.LiveDifficultyID,
			baseTime:         baseDate,
		})
	}
	sort.Ints(groupIDs)

	result := make([]specialLiveRotationSchedule, 0, len(rows))
	for _, groupID := range groupIDs {
		liveList := grouped[groupID]
		if len(liveList) == 0 {
			continue
		}

		periodDays := int64(1)
		if len(liveList) >= 2 {
			diff := int64(liveList[1].baseTime.Sub(liveList[0].baseTime) / (24 * time.Hour))
			if diff > 0 {
				periodDays = diff
			}
		}
		if periodDays != 1 {
			continue
		}

		for _, live := range liveList {
			result = append(result, specialLiveRotationSchedule{
				LiveDifficultyID: live.liveDifficultyID,
				StartTime:        time.Unix(0, 0).In(cst),
				EndTime:          time.Unix(2147483647, 0).In(cst),
			})
		}
	}

	return result, nil
}

func listAvailableTrainingLiveDifficultyIDs(ss *session.Session) ([]int, error) {
	specialRows, err := listAvailableSpecialLiveRows(ss)
	if err != nil {
		return nil, err
	}

	result := make([]int, 0)
	for _, row := range specialRows {
		if row.Difficulty <= 5 || row.AcFlag != 0 || row.ExcludeClearCountFlag != 0 {
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
		t, err := time.ParseInLocation(layout, value, cst)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid special live rotation base_date: %s", value)
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func dayIndexInCST(t time.Time) int64 {
	return startOfDay(t.In(cst)).Unix() / 86400
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
