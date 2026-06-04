package live

import (
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"honoka-chan/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func reward(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer func() {
		ss.ClearLiveInProgress()
		ss.Finalize()
	}()

	playRewardReq := liveschema.RewardReq{}
	err := honokautils.ParseRequestData(ctx, &playRewardReq)
	if ss.CheckErr(err) {
		return
	}
	// fmt.Println(ctx.MustGet("request_data").(string))

	difficultyID := playRewardReq.LiveDifficultyID
	_, liveInfo := ss.GetLiveInfo(difficultyID)
	_, progress := ss.GetLiveInProgress() // TODO: 添加返回指定状态码的方法
	_, deckInfo := ss.GetDeckInfo(progress.DeckID)
	deckInfo.LiveDifficultyID = difficultyID

	unitsList := []liveschema.PlayRewardUnitList{}
	err = ss.UserEng.Table("user_deck_unit").
		Join("LEFT", "user_deck", "user_deck_unit.user_deck_id = user_deck.id").
		Where("user_deck.user_id = ? AND user_deck.deck_id = ?", ss.UserID, progress.DeckID).
		Find(&unitsList)
	if ss.CheckErr(err) {
		return
	}

	totalScore := playRewardReq.ScoreSmile + playRewardReq.ScoreCool + playRewardReq.ScoreCute
	playResp := liveschema.RewardResp{
		ResponseData: liveschema.RewardData{
			LiveInfo: []liveschema.RewardLiveInfo{
				{
					LiveDifficultyID: difficultyID,
					IsRandom:         false,
					AcFlag:           liveInfo.AcFlag,
					SwingFlag:        liveInfo.SwingFlag,
				},
			},
			TotalLove:   0,
			IsHighScore: true,
			HiScore:     totalScore,
			BaseRewardInfo: liveschema.BaseRewardInfo{
				PlayerExp: 0,
				PlayerExpUnitMax: liveschema.PlayerExpUnitMax{
					Before: 0,
					After:  0,
				},
				PlayerExpFriendMax: liveschema.PlayerExpFriendMax{
					Before: 99,
					After:  99,
				},
				PlayerExpLpMax: liveschema.PlayerExpLpMax{
					Before: ss.UserPref.EffectiveEnergyMax(),
					After:  ss.UserPref.EffectiveEnergyMax(),
				},
				GameCoin:              0,
				GameCoinRewardBoxFlag: false,
				SocialPoint:           0,
			},
			RewardUnitList: liveschema.RewardUnitList{
				LiveClear: []liveschema.LiveClear{},
				LiveRank:  []liveschema.LiveRank{},
				LiveCombo: []any{},
			},
			UnlockedSubscenarioIds:       []any{},
			UnlockedMultiUnitScenarioIds: []any{},
			EffortPoint:                  []liveschema.EffortPoint{},
			IsEffortPointVisible:         false,
			LimitedEffortBox:             []any{},
			UnitList:                     unitsList,
			BeforeUserInfo:               ss.GetUserInfo(),
			AfterUserInfo:                ss.GetUserInfo(),
			NextLevelInfo: []liveschema.NextLevelInfo{
				{
					Level:   ss.UserPref.UserLevel,
					FromExp: ss.UserPref.UserExp,
				},
			},
			GoalAccompInfo: liveschema.GoalAccompInfo{
				AchievedIds: []any{},
				Rewards:     []any{},
			},
			SpecialRewardInfo:    []any{},
			EventInfo:            []any{},
			DailyRewardInfo:      []any{},
			CanSendFriendRequest: false,
			UsingBuffInfo:        []any{},
			ClassSystem: liveschema.ClassSystem{
				RankInfo: liveschema.RewardRankInfo{
					BeforeClassRankID: 10,
					AfterClassRankID:  10,
					RankUpDate:        "2020-02-12 11:57:15",
				},
				CompleteFlag: false,
				IsOpened:     true,
				IsVisible:    true,
			},
			AccomplishedAchievementList:  []liveschema.AccomplishedAchievementList{},
			UnaccomplishedAchievementCnt: 0,
			AddedAchievementList:         []any{},
			MuseumInfo:                   liveschema.Museum{},
			UnitSupportList:              []liveschema.RewardUnitSupportList{},
			ServerTimestamp:              time.Now().Unix(),
			PresentCnt:                   0,
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	if playRewardReq.MaxCombo > liveInfo.SRankCombo {
		playResp.ResponseData.ComboRank = 1
	} else if playRewardReq.MaxCombo > liveInfo.ARankCombo {
		playResp.ResponseData.ComboRank = 2
	} else if playRewardReq.MaxCombo > liveInfo.BRankCombo {
		playResp.ResponseData.ComboRank = 3
	} else if playRewardReq.MaxCombo > liveInfo.CRankCombo {
		playResp.ResponseData.ComboRank = 4
	} else {
		playResp.ResponseData.ComboRank = 5
	}

	if totalScore > liveInfo.SRankScore {
		playResp.ResponseData.Rank = 1
	} else if totalScore > liveInfo.ARankScore {
		playResp.ResponseData.Rank = 2
	} else if totalScore > liveInfo.BRankScore {
		playResp.ResponseData.Rank = 3
	} else if totalScore > liveInfo.CRankScore {
		playResp.ResponseData.Rank = 4
	} else {
		playResp.ResponseData.Rank = 5
	}

	// 判断歌曲记录是否需要保存或者更新
	liveSetting := playRewardReq.PreciseScoreLog.LiveSetting
	liveSetting.BackgroundID = ss.UserPref.BackgroundID
	if playRewardReq.PreciseScoreLog.IsLogOn && liveSetting.PreciseScoreAutoUpdateFlag {
		newLiveRecord := usermodel.UserLiveRecord{
			LiveDifficultyID: difficultyID,
			DeckID:           progress.DeckID,
			TotalScore:       totalScore,
			PerfectCnt:       playRewardReq.PerfectCnt,
			GreatCnt:         playRewardReq.GreatCnt,
			GoodCnt:          playRewardReq.GoodCnt,
			BadCnt:           playRewardReq.BadCnt,
			MissCnt:          playRewardReq.MissCnt,
			MaxCombo:         playRewardReq.MaxCombo,
			IsSkillOn:        playRewardReq.PreciseScoreLog.IsSkillOn,
			LiveInfoJSON:     utils.ToJSON(liveInfo),
			PreciseListJSON:  utils.ToJSON(playRewardReq.PreciseScoreLog.PreciseList),
			DeckInfoJSON:     utils.ToJSON(deckInfo),
			TapAdjust:        playRewardReq.PreciseScoreLog.TapAdjust,
			LiveSettingJSON:  utils.ToJSON(liveSetting),
			CanReplay:        true,
			TriggerLogJSON:   utils.ToJSON(playRewardReq.PreciseScoreLog.TriggerLog),
			UserID:           ss.UserID,
			UpdateDate:       time.Now().Format("2006-01-02 15:04:05"),
		}
		ss.UpdateUserLiveRecord(&newLiveRecord, liveSetting.PreciseScoreUpdateType)
	}

	ss.Respond(playResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/live/reward", middleware.Common, reward)
}
