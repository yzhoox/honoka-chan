package live

import (
	"encoding/json"
	"honoka-chan/internal/constant"
	"honoka-chan/internal/middleware"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/router"
	liveschema "honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	"honoka-chan/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func preciseScore(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	playScoreReq := liveschema.PlayScoreReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &playScoreReq)
	if ss.CheckErr(err) {
		return
	}
	// fmt.Println(ctx.MustGet("request_data").(string))

	difficultyID, _ := strconv.Atoi(playScoreReq.LiveDifficultyID)

	// 歌曲类型: normal / special
	// sqlite3 不支持 FULL OUTER JOIN 所以这里使用 UNION ALL
	var liveSetting struct {
		NotesSettingAsset string `xorm:"notes_setting_asset"`
		ARankScore        int    `xorm:"a_rank_score"`
		BRankScore        int    `xorm:"b_rank_score"`
		CRankScore        int    `xorm:"c_rank_score"`
		SRankScore        int    `xorm:"s_rank_score"`
		AcFlag            int    `xorm:"ac_flag"`
		SwingFlag         int    `xorm:"swing_flag"`
	}
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
	_, err = ss.MainEng.SQL(sql, difficultyID, difficultyID).Get(&liveSetting)
	if ss.CheckErr(err) {
		return
	}
	// fmt.Println("liveSetting", liveSetting)

	notesList := []liveschema.NotesList{}
	noteData := utils.ReadAllText("./assets/serverdata/beatmaps/" + liveSetting.NotesSettingAsset)
	err = json.Unmarshal([]byte(noteData), &notesList)
	if ss.CheckErr(err) {
		return
	}

	ranks := []liveschema.RankInfo{}
	ranks = append(ranks, liveschema.RankInfo{
		Rank:    5,
		RankMin: 0,
		RankMax: liveSetting.CRankScore,
	}, liveschema.RankInfo{
		Rank:    4,
		RankMin: liveSetting.CRankScore + 1,
		RankMax: liveSetting.BRankScore,
	}, liveschema.RankInfo{
		Rank:    3,
		RankMin: liveSetting.BRankScore + 1,
		RankMax: liveSetting.ARankScore,
	}, liveschema.RankInfo{
		Rank:    2,
		RankMin: liveSetting.ARankScore + 1,
		RankMax: liveSetting.SRankScore,
	}, liveschema.RankInfo{
		Rank:    1,
		RankMin: liveSetting.SRankScore + 1,
		RankMax: 0,
	})

	// 检查是否正在进行 Live
	progress, _ := ss.GetLiveInProgress()

	// 检查是否有 Live 记录
	liveRecord := ss.GetUserLiveRecord(difficultyID)

	var skillOn liveschema.Skill
	var skillOff liveschema.Skill
	var playResp liveschema.PreciseScoreResp

	// 正在进行 Live
	if progress {
		// 返回默认
		playResp = liveschema.PreciseScoreResp{
			ResponseData: liveschema.PreciseScoreData{
				On: liveschema.Skill{
					HasRecord: false,
					LiveInfo: liveschema.LiveInfo{
						LiveDifficultyID: difficultyID,
						IsRandom:         false,
						AcFlag:           liveSetting.AcFlag,
						SwingFlag:        liveSetting.SwingFlag,
						NotesList:        notesList,
					},
				},
				Off: liveschema.Skill{
					HasRecord: false,
					LiveInfo: liveschema.LiveInfo{
						LiveDifficultyID: difficultyID,
						IsRandom:         false,
						AcFlag:           liveSetting.AcFlag,
						SwingFlag:        liveSetting.SwingFlag,
						NotesList:        notesList,
					},
				},
				RankInfo:          ranks,
				CanActivateEffect: true,
				ServerTimestamp:   time.Now().Unix(),
			},
			ReleaseInfo: []any{},
			StatusCode:  200,
		}
	} else {
		// 如果有 Live 记录
		if len(liveRecord) > 0 {
			skillOn = liveschema.Skill{
				HasRecord:   false,
				RandomSeed:  nil,
				MaxCombo:    nil,
				UpdateDate:  nil,
				PreciseList: nil,
				DeckInfo:    nil,
				TapAdjust:   nil,
				CanReplay:   false,
			}
			skillOff = skillOn

			// 已完成的 Live 判断技能开关情况
			for _, record := range liveRecord {
				// LiveInfo
				var liveInfo liveschema.LiveInfo
				err = json.Unmarshal([]byte(record.LiveInfoJSON), &liveInfo)
				if ss.CheckErr(err) {
					return
				}
				liveInfo.NotesList = notesList
				skillOn.LiveInfo = liveInfo
				skillOff.LiveInfo = liveInfo

				// PreciseList
				var preciseList []liveschema.PreciseList
				err = json.Unmarshal([]byte(record.PreciseListJSON), &preciseList)
				if ss.CheckErr(err) {
					return
				}

				// DeckInfo
				var deckInfo usermodel.DeckInfo
				err = json.Unmarshal([]byte(record.DeckInfoJSON), &deckInfo)
				if ss.CheckErr(err) {
					return
				}

				// LiveSetting
				var liveSetting liveschema.LiveSetting
				err = json.Unmarshal([]byte(record.LiveSettingJSON), &liveSetting)
				if ss.CheckErr(err) {
					return
				}

				// TriggerLog
				var triggerLog []liveschema.TriggerLog
				err = json.Unmarshal([]byte(record.TriggerLogJSON), &triggerLog)
				if ss.CheckErr(err) {
					return
				}

				if record.IsSkillOn {
					// 技能开
					skillOn.HasRecord = true
					skillOn.RandomSeed = time.Now().Unix() // TODO: 从 /live/play 的 Timestamp 字段获取
					skillOn.MaxCombo = record.MaxCombo
					skillOn.UpdateDate = record.UpdateDate
					skillOn.PreciseList = preciseList
					skillOn.DeckInfo = deckInfo
					skillOn.LiveSetting = liveSetting
					skillOn.TriggerLog = triggerLog
					skillOn.TapAdjust = record.TapAdjust
					skillOn.CanReplay = record.CanReplay
				} else {
					// 技能关
					skillOff.HasRecord = true
					skillOff.RandomSeed = time.Now().Unix() // TODO: 从 /live/play 的 Timestamp 字段获取
					skillOff.MaxCombo = record.MaxCombo
					skillOff.UpdateDate = record.UpdateDate
					skillOff.PreciseList = preciseList
					skillOff.DeckInfo = deckInfo
					skillOff.LiveSetting = liveSetting
					skillOff.TriggerLog = triggerLog
					skillOff.TapAdjust = record.TapAdjust
					skillOff.CanReplay = record.CanReplay
				}
			}

			playResp = liveschema.PreciseScoreResp{
				ResponseData: liveschema.PreciseScoreData{
					On:                skillOn,
					Off:               skillOff,
					RankInfo:          ranks,
					CanActivateEffect: true,
					ServerTimestamp:   time.Now().Unix(),
				},
				ReleaseInfo: []any{},
				StatusCode:  200,
			}
		} else {
			playResp = liveschema.PreciseScoreResp{
				ResponseData: map[string]constant.ErrorCode{
					"error_code": constant.ErrorCodeLivePreciseListNotFound,
				},
				ReleaseInfo: []any{},
				StatusCode:  600,
			}
		}
	}

	ss.Respond(playResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/live/preciseScore", middleware.Common, preciseScore)
}
