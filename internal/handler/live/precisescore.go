package live

import (
	"encoding/json"
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/live"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func preciseScore(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	playScoreReq := live.PlayScoreReq{}
	err := json.Unmarshal([]byte(ctx.MustGet("request_data").(string)), &playScoreReq)
	if ss.CheckErr(err) {
		return
	}

	tDifficultyId := playScoreReq.LiveDifficultyID
	difficultyId, err := strconv.Atoi(tDifficultyId)
	if ss.CheckErr(err) {
		return
	}

	// Song type: normal / special
	// sqlite3 doesn't support FULL OUTER JOIN so use UNION ALL here.
	sql := `SELECT notes_setting_asset,c_rank_score,b_rank_score,a_rank_score,s_rank_score,ac_flag,swing_flag FROM live_setting_m WHERE live_setting_id IN (SELECT live_setting_id FROM normal_live_m WHERE live_difficulty_id = ? UNION ALL SELECT live_setting_id FROM special_live_m WHERE live_difficulty_id = ?)`
	var notes_setting_asset string
	var c_rank_score, b_rank_score, a_rank_score, s_rank_score, ac_flag, swing_flag int
	err = ss.MainEng.DB().QueryRow(sql, difficultyId, difficultyId).Scan(&notes_setting_asset, &c_rank_score, &b_rank_score, &a_rank_score, &s_rank_score, &ac_flag, &swing_flag)
	if ss.CheckErr(err) {
		return
	}

	// fmt.Println(notes_setting_asset)
	// fmt.Println(c_rank_score, b_rank_score, a_rank_score, s_rank_score)

	notes := []live.NotesList{}
	// fmt.Println("./assets/serverdata/beatmaps/" + notes_setting_asset)
	notes_list := honokautils.ReadAllText("./assets/serverdata/beatmaps/" + notes_setting_asset)
	err = json.Unmarshal([]byte(notes_list), &notes)
	if ss.CheckErr(err) {
		return
	}

	ranks := []live.RankInfo{}
	ranks = append(ranks, live.RankInfo{
		Rank:    5,
		RankMin: 0,
		RankMax: c_rank_score,
	}, live.RankInfo{
		Rank:    4,
		RankMin: c_rank_score + 1,
		RankMax: b_rank_score,
	}, live.RankInfo{
		Rank:    3,
		RankMin: b_rank_score + 1,
		RankMax: a_rank_score,
	}, live.RankInfo{
		Rank:    2,
		RankMin: a_rank_score + 1,
		RankMax: s_rank_score,
	}, live.RankInfo{
		Rank:    1,
		RankMin: s_rank_score + 1,
		RankMax: 0,
	})

	playResp := live.PreciseScoreResp{
		ResponseData: live.PreciseScoreData{
			On: live.On{
				HasRecord: false,
				LiveInfo: live.LiveInfo{
					LiveDifficultyID: difficultyId,
					IsRandom:         false,
					AcFlag:           ac_flag,
					SwingFlag:        swing_flag,
					NotesList:        notes,
				},
			},
			Off: live.Off{
				HasRecord: false,
				LiveInfo: live.LiveInfo{
					LiveDifficultyID: difficultyId,
					IsRandom:         false,
					AcFlag:           ac_flag,
					SwingFlag:        swing_flag,
					NotesList:        notes,
				},
			},
			RankInfo:          ranks,
			CanActivateEffect: true,
			ServerTimestamp:   time.Now().Unix(),
		},
		ReleaseInfo: []any{},
		StatusCode:  200,
	}

	ss.Respond(playResp)
}

func init() {
	router.AddHandler("main.php", "POST", "/live/preciseScore", middleware.Common, preciseScore)
}
