package live

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/db"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func reward(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	playRewardReq := live.RewardReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &playRewardReq)
	if ss.CheckErr(err) {
		return
	}

	difficultyId := playRewardReq.LiveDifficultyID

	// Song type: normal / special
	// sqlite3 doesn't support FULL OUTER JOIN so use UNION ALL here.
	sql := `SELECT c_rank_score,b_rank_score,a_rank_score,s_rank_score,c_rank_combo,b_rank_combo,a_rank_combo,s_rank_combo,ac_flag,swing_flag FROM live_setting_m WHERE live_setting_id IN (SELECT live_setting_id FROM normal_live_m WHERE live_difficulty_id = ? UNION ALL SELECT live_setting_id FROM special_live_m WHERE live_difficulty_id = ?)`
	var c_rank_score, b_rank_score, a_rank_score, s_rank_score, c_rank_combo, b_rank_combo, a_rank_combo, s_rank_combo, ac_flag, swing_flag int
	err = ss.MainEng.DB().QueryRow(sql, difficultyId, difficultyId).Scan(&c_rank_score, &b_rank_score, &a_rank_score, &s_rank_score, &c_rank_combo, &b_rank_combo, &a_rank_combo, &s_rank_combo, &ac_flag, &swing_flag)
	if ss.CheckErr(err) {
		return
	}

	key := "live_deck_" + strconv.Itoa(ss.UserID)
	deckId, err := db.Ldb.Get([]byte(key))
	if ss.CheckErr(err) {
		return
	}
	unitsList := []live.PlayRewardUnitList{}
	err = ss.UserEng.Table("user_deck_unit").Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
		Where("user_id = ? AND deck_id = ?", ss.UserID, string(deckId)).Find(&unitsList)
	if ss.CheckErr(err) {
		return
	}

	pref := user.UserPref{}
	_, err = ss.UserEng.Table("user_pref").Where("user_id = ?", ss.UserID).Get(&pref)
	if ss.CheckErr(err) {
		return
	}

	totalScore := playRewardReq.ScoreSmile + playRewardReq.ScoreCool + playRewardReq.ScoreCute
	playResp := live.RewardResp{
		ResponseData: live.RewardData{
			LiveInfo: []live.RewardLiveInfo{
				{
					LiveDifficultyID: difficultyId,
					IsRandom:         false,
					AcFlag:           ac_flag,
					SwingFlag:        swing_flag,
				},
			},
			TotalLove:   0,
			IsHighScore: true,
			HiScore:     totalScore,
			BaseRewardInfo: live.BaseRewardInfo{
				PlayerExp: 0,
				PlayerExpUnitMax: live.PlayerExpUnitMax{
					Before: 0,
					After:  0,
				},
				PlayerExpFriendMax: live.PlayerExpFriendMax{
					Before: 99,
					After:  99,
				},
				PlayerExpLpMax: live.PlayerExpLpMax{
					Before: 417,
					After:  417,
				},
				GameCoin:              0,
				GameCoinRewardBoxFlag: false,
				SocialPoint:           0,
			},
			RewardUnitList: live.RewardUnitList{
				LiveClear: []live.LiveClear{},
				LiveRank:  []live.LiveRank{},
				LiveCombo: []any{},
			},
			UnlockedSubscenarioIds:       []any{},
			UnlockedMultiUnitScenarioIds: []any{},
			EffortPoint:                  []live.EffortPoint{},
			IsEffortPointVisible:         false,
			LimitedEffortBox:             []any{},
			UnitList:                     unitsList,
			BeforeUserInfo: live.BeforeUserInfo{
				Level:                          pref.UserLevel,
				Exp:                            1089696,
				PreviousExp:                    0,
				NextExp:                        1207185,
				GameCoin:                       999999999,
				SnsCoin:                        10000,
				FreeSnsCoin:                    50000,
				PaidSnsCoin:                    50000,
				SocialPoint:                    1438165,
				UnitMax:                        5000,
				WaitingUnitMax:                 1000,
				CurrentEnergy:                  417,
				EnergyMax:                      417,
				TrainingEnergy:                 9,
				TrainingEnergyMax:              10,
				EnergyFullTime:                 "2023-03-20 01:28:55",
				LicenseLiveEnergyRecoverlyTime: 60,
				FriendMax:                      99,
				TutorialState:                  -1,
				OverMaxEnergy:                  417,
				UnlockRandomLiveMuse:           1,
				UnlockRandomLiveAqours:         1,
			},
			AfterUserInfo: live.AfterUserInfo{
				Level:                          pref.UserLevel,
				Exp:                            1089696,
				PreviousExp:                    0,
				NextExp:                        1207185,
				GameCoin:                       999999999,
				SnsCoin:                        100000,
				FreeSnsCoin:                    50000,
				PaidSnsCoin:                    50000,
				SocialPoint:                    1438375,
				UnitMax:                        5000,
				WaitingUnitMax:                 1000,
				CurrentEnergy:                  417,
				EnergyMax:                      417,
				TrainingEnergy:                 9,
				TrainingEnergyMax:              10,
				EnergyFullTime:                 "2023-03-20 01:28:55",
				LicenseLiveEnergyRecoverlyTime: 60,
				FriendMax:                      99,
				TutorialState:                  -1,
				OverMaxEnergy:                  417,
				UnlockRandomLiveMuse:           1,
				UnlockRandomLiveAqours:         1,
			},
			NextLevelInfo: []live.NextLevelInfo{
				{
					Level:   pref.UserLevel,
					FromExp: 1089696,
				},
			},
			GoalAccompInfo: live.GoalAccompInfo{
				AchievedIds: []any{},
				Rewards:     []any{},
			},
			SpecialRewardInfo:    []any{},
			EventInfo:            []any{},
			DailyRewardInfo:      []any{},
			CanSendFriendRequest: false,
			UsingBuffInfo:        []any{},
			ClassSystem: live.ClassSystem{
				RankInfo: live.RewardRankInfo{
					BeforeClassRankID: 10,
					AfterClassRankID:  10,
					RankUpDate:        "2020-02-12 11:57:15",
				},
				CompleteFlag: false,
				IsOpened:     true,
				IsVisible:    true,
			},
			AccomplishedAchievementList:  []live.AccomplishedAchievementList{},
			UnaccomplishedAchievementCnt: 0,
			AddedAchievementList:         []any{},
			MuseumInfo:                   live.Museum{},
			UnitSupportList:              []live.RewardUnitSupportList{},
			ServerTimestamp:              time.Now().Unix(),
			PresentCnt:                   0,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	if playRewardReq.MaxCombo > s_rank_combo {
		playResp.ResponseData.ComboRank = 1
	} else if playRewardReq.MaxCombo > a_rank_combo {
		playResp.ResponseData.ComboRank = 2
	} else if playRewardReq.MaxCombo > b_rank_combo {
		playResp.ResponseData.ComboRank = 3
	} else if playRewardReq.MaxCombo > c_rank_combo {
		playResp.ResponseData.ComboRank = 4
	} else {
		playResp.ResponseData.ComboRank = 5
	}

	if totalScore > s_rank_score {
		playResp.ResponseData.Rank = 1
	} else if totalScore > a_rank_score {
		playResp.ResponseData.Rank = 2
	} else if totalScore > b_rank_score {
		playResp.ResponseData.Rank = 3
	} else if totalScore > c_rank_score {
		playResp.ResponseData.Rank = 4
	} else {
		playResp.ResponseData.Rank = 5
	}

	ss.Respond(playResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/live/reward", middleware.Common, reward)
}
