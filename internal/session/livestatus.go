package session

import (
	"honoka-chan/internal/constant"
	usermodel "honoka-chan/internal/model/user"
	"sort"
	"time"
)

type liveGoalRewardRow struct {
	LiveGoalRewardID int                   `xorm:"live_goal_reward_id"`
	LiveDifficultyID int                   `xorm:"live_difficulty_id"`
	LiveGoalType     constant.LiveGoalType `xorm:"live_goal_type"`
	Rank             int                   `xorm:"rank"`
	AddType          int                   `xorm:"add_type"`
	ItemID           int                   `xorm:"item_id"`
	ItemCategoryID   int                   `xorm:"item_category_id"`
	Amount           int                   `xorm:"amount"`
}

type LiveStatusSnapshot struct {
	Status             int
	HiScore            int
	HiComboCount       int
	ClearCnt           int
	AchievedGoalIDList []int
}

func liveRank(value int, cRank int, bRank int, aRank int, sRank int) int {
	switch {
	case value >= sRank:
		return 1
	case value >= aRank:
		return 2
	case value >= bRank:
		return 3
	case value >= cRank:
		return 4
	default:
		return 5
	}
}

func (ss *Session) GetUserLiveStatus(liveDifficultyID int) (bool, *usermodel.UserLiveStatus, error) {
	status := usermodel.UserLiveStatus{}
	has, err := ss.UserEng.Table(new(usermodel.UserLiveStatus)).
		Where("user_id = ? AND live_difficulty_id = ?", ss.UserID, liveDifficultyID).
		Get(&status)
	if err != nil {
		return false, nil, err
	}
	return has, &status, nil
}

func (ss *Session) GetUserLiveStatusMap(liveDifficultyIDs []int) (map[int]usermodel.UserLiveStatus, error) {
	result := map[int]usermodel.UserLiveStatus{}
	if len(liveDifficultyIDs) == 0 {
		return result, nil
	}

	rows := []usermodel.UserLiveStatus{}
	err := ss.UserEng.Table(new(usermodel.UserLiveStatus)).
		In("live_difficulty_id", liveDifficultyIDs).
		Where("user_id = ?", ss.UserID).
		Find(&rows)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.LiveDifficultyID] = row
	}
	return result, nil
}

func (ss *Session) GetUserLiveGoalRewardIDMap(liveDifficultyIDs []int) (map[int][]int, error) {
	result := map[int][]int{}
	if len(liveDifficultyIDs) == 0 {
		return result, nil
	}

	rows := []usermodel.UserLiveGoal{}
	err := ss.UserEng.Table(new(usermodel.UserLiveGoal)).
		In("live_difficulty_id", liveDifficultyIDs).
		Where("user_id = ?", ss.UserID).
		Where("live_goal_reward_id > 0").
		OrderBy("live_goal_reward_id ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.LiveDifficultyID] = append(result[row.LiveDifficultyID], row.LiveGoalRewardID)
	}
	return result, nil
}

func (ss *Session) BuildLiveStatusSnapshotMap(liveDifficultyIDs []int) (map[int]LiveStatusSnapshot, error) {
	statusMap, err := ss.GetUserLiveStatusMap(liveDifficultyIDs)
	if err != nil {
		return nil, err
	}
	goalMap, err := ss.GetUserLiveGoalRewardIDMap(liveDifficultyIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[int]LiveStatusSnapshot, len(liveDifficultyIDs))
	for _, liveDifficultyID := range liveDifficultyIDs {
		status := LiveStatusSnapshot{
			Status:             1,
			HiScore:            0,
			HiComboCount:       0,
			ClearCnt:           0,
			AchievedGoalIDList: []int{},
		}

		if row, ok := statusMap[liveDifficultyID]; ok {
			status.HiScore = row.HiScore
			status.HiComboCount = row.HiComboCount
			status.ClearCnt = row.ClearCnt
			if row.ClearCnt > 0 {
				status.Status = 2
			}
		}

		if ids := goalMap[liveDifficultyID]; len(ids) > 0 {
			status.AchievedGoalIDList = ids
		}

		result[liveDifficultyID] = status
	}
	return result, nil
}

func (ss *Session) UpdateUserLiveStatus(liveDifficultyID int, totalScore int, maxCombo int) (before usermodel.UserLiveStatus, after usermodel.UserLiveStatus, err error) {
	now := time.Now().Unix()
	has, current, err := ss.GetUserLiveStatus(liveDifficultyID)
	if err != nil {
		return before, after, err
	}
	if current == nil {
		current = &usermodel.UserLiveStatus{}
	}

	before = *current
	if !has {
		after = usermodel.UserLiveStatus{
			UserID:           ss.UserID,
			LiveDifficultyID: liveDifficultyID,
			HiScore:          totalScore,
			HiComboCount:     maxCombo,
			ClearCnt:         1,
			InsertDate:       now,
			UpdateDate:       now,
		}
		_, err = ss.UserEng.Insert(&after)
		return before, after, err
	}

	after = before
	if totalScore > after.HiScore {
		after.HiScore = totalScore
	}
	if maxCombo > after.HiComboCount {
		after.HiComboCount = maxCombo
	}
	after.ClearCnt++
	after.UpdateDate = now

	_, err = ss.UserEng.Table(new(usermodel.UserLiveStatus)).
		Where("user_id = ? AND live_difficulty_id = ?", ss.UserID, liveDifficultyID).
		Cols("hi_score", "hi_combo_count", "clear_cnt", "update_date").
		Update(&after)
	return before, after, err
}

func (ss *Session) GetAchievedLiveGoalRewardRows(liveDifficultyID int, hiScore int, hiComboCount int, clearCnt int) ([]liveGoalRewardRow, error) {
	has, liveInfo := ss.GetLiveInfo(liveDifficultyID)
	if !has || liveInfo == nil {
		return []liveGoalRewardRow{}, nil
	}

	rows := []liveGoalRewardRow{}
	err := ss.MainEng.Table("live_goal_reward_m").
		Where("live_difficulty_id = ?", liveDifficultyID).
		OrderBy("live_goal_type ASC, rank ASC, live_goal_reward_id ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	scoreRank := liveRank(hiScore, liveInfo.CRankScore, liveInfo.BRankScore, liveInfo.ARankScore, liveInfo.SRankScore)
	comboRank := liveRank(hiComboCount, liveInfo.CRankCombo, liveInfo.BRankCombo, liveInfo.ARankCombo, liveInfo.SRankCombo)
	clearRank := liveRank(clearCnt, liveInfo.CRankComplete, liveInfo.BRankComplete, liveInfo.ARankComplete, liveInfo.SRankComplete)

	result := make([]liveGoalRewardRow, 0, len(rows))
	for _, row := range rows {
		switch row.LiveGoalType {
		case constant.LiveGoalTypeScore:
			if scoreRank <= row.Rank {
				result = append(result, row)
			}
		case constant.LiveGoalTypeCombo:
			if comboRank <= row.Rank {
				result = append(result, row)
			}
		case constant.LiveGoalTypeClear:
			if clearRank <= row.Rank {
				result = append(result, row)
			}
		}
	}
	return result, nil
}

func (ss *Session) SyncUserLiveGoals(liveDifficultyID int, hiScore int, hiComboCount int, clearCnt int) ([]int, error) {
	achievedRows, err := ss.GetAchievedLiveGoalRewardRows(liveDifficultyID, hiScore, hiComboCount, clearCnt)
	if err != nil {
		return nil, err
	}
	if len(achievedRows) == 0 {
		return []int{}, nil
	}

	existingRows := []usermodel.UserLiveGoal{}
	err = ss.UserEng.Table(new(usermodel.UserLiveGoal)).
		Where("user_id = ? AND live_difficulty_id = ?", ss.UserID, liveDifficultyID).
		Where("live_goal_reward_id > 0").
		Find(&existingRows)
	if err != nil {
		return nil, err
	}

	existing := make(map[int]struct{}, len(existingRows))
	for _, row := range existingRows {
		existing[row.LiveGoalRewardID] = struct{}{}
	}

	now := time.Now().Unix()
	newlyAchieved := make([]int, 0)
	for _, row := range achievedRows {
		if _, ok := existing[row.LiveGoalRewardID]; ok {
			continue
		}

		_, err = ss.UserEng.Insert(&usermodel.UserLiveGoal{
			UserID:           ss.UserID,
			LiveDifficultyID: liveDifficultyID,
			LiveGoalRewardID: row.LiveGoalRewardID,
			GoalType:         row.LiveGoalType,
			Rank:             row.Rank,
			CompletedAt:      now,
		})
		if err != nil {
			return nil, err
		}

		newlyAchieved = append(newlyAchieved, row.LiveGoalRewardID)
	}

	sort.Ints(newlyAchieved)
	return newlyAchieved, nil
}

func (ss *Session) GetLiveGoalRewardRowsByIDs(goalRewardIDs []int) ([]liveGoalRewardRow, error) {
	if len(goalRewardIDs) == 0 {
		return []liveGoalRewardRow{}, nil
	}

	rows := []liveGoalRewardRow{}
	err := ss.MainEng.Table("live_goal_reward_m").
		In("live_goal_reward_id", goalRewardIDs).
		OrderBy("live_goal_reward_id ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
