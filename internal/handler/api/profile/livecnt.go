package profile

import (
	usermodel "honoka-chan/internal/model/user"
	profileapischema "honoka-chan/internal/schema/api/profile"
	"honoka-chan/internal/session"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type profileLiveDifficultyRow struct {
	LiveDifficultyID int `xorm:"live_difficulty_id"`
	Difficulty       int `xorm:"difficulty"`
}

func liveCnt(ctx *gin.Context, targetUserID int) (res any, err error) {
	ss := session.Get(ctx)

	targetPref, err := getTargetUserPref(ss, targetUserID)
	if err != nil {
		return nil, err
	}

	difficultyByLiveID := map[int]int{}
	loadDifficultyRows := func(table string) error {
		rows := []profileLiveDifficultyRow{}
		err := ss.MainEng.Table(table).Alias("live").
			Join("LEFT", "live_setting_m setting", "live.live_setting_id = setting.live_setting_id").
			Select("live.live_difficulty_id, setting.difficulty").
			Find(&rows)
		if err != nil {
			return err
		}
		for _, row := range rows {
			difficultyByLiveID[row.LiveDifficultyID] = row.Difficulty
		}
		return nil
	}
	if err := loadDifficultyRows("normal_live_m"); err != nil {
		return nil, err
	}
	if err := loadDifficultyRows("special_live_m"); err != nil {
		return nil, err
	}

	statusRows := []usermodel.UserLiveStatus{}
	err = ss.UserEng.Table(new(usermodel.UserLiveStatus)).
		Where("user_id = ?", targetPref.UserID).
		Where("clear_cnt > 0").
		Cols("live_difficulty_id", "clear_cnt").
		Find(&statusRows)
	if err != nil {
		return nil, err
	}

	difficultyCounts := map[int]int{}
	for _, row := range statusRows {
		difficulty, ok := difficultyByLiveID[row.LiveDifficultyID]
		if !ok {
			continue
		}
		if difficulty == 5 {
			continue
		}
		difficultyCounts[difficulty]++
	}

	orderedDifficulties := []int{1, 2, 3, 4, 6}
	result := make([]profileapischema.LiveCntData, 0, len(orderedDifficulties))
	for _, difficulty := range orderedDifficulties {
		result = append(result, profileapischema.LiveCntData{
			Difficulty: difficulty,
			ClearCnt:   difficultyCounts[difficulty],
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Difficulty < result[j].Difficulty
	})

	res = profileapischema.LiveCntResp{
		Result:     result,
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
