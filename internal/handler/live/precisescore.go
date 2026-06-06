package live

import (
	"encoding/json"
	"honoka-chan/internal/constant"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	commonschema "honoka-chan/internal/schema/common"
	liveschema "honoka-chan/internal/schema/live"
	liverecordschema "honoka-chan/internal/schema/liverecord"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func preciseScore(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	playScoreReq := liveschema.PlayScoreReq{}
	err := honokautils.ParseRequestData(ctx, &playScoreReq)
	if ss.CheckErr(err) {
		return
	}
	// fmt.Println(ctx.MustGet("request_data").(string))

	difficultyID, _ := strconv.Atoi(playScoreReq.LiveDifficultyID)

	liveSetting, err := loadLiveSetting(ss, difficultyID)
	if ss.CheckErr(err) {
		return
	}

	notesList, err := loadLiveNotes(liveSetting.NotesSettingAsset)
	if ss.CheckErr(err) {
		return
	}

	ranks := buildRankInfo(liveSetting)

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
				var deckInfo liverecordschema.DeckInfo
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
				ResponseData: commonschema.ErrorData{
					ErrorCode: constant.ErrorCodeLivePreciseListNotFound,
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
