package rliveschema

import liveschema "honoka-chan/internal/schema/live"

type RewardReq struct {
	Module           string                     `json:"module"`
	Action           string                     `json:"action"`
	GoodCnt          int                        `json:"good_cnt"`
	MissCnt          int                        `json:"miss_cnt"`
	IsTraining       bool                       `json:"is_training"`
	GreatCnt         int                        `json:"great_cnt"`
	CommandNum       string                     `json:"commandNum"`
	LoveCnt          int                        `json:"love_cnt"`
	RemainHp         int                        `json:"remain_hp"`
	MaxCombo         int                        `json:"max_combo"`
	ScoreSmile       int                        `json:"score_smile"`
	PerfectCnt       int                        `json:"perfect_cnt"`
	BadCnt           int                        `json:"bad_cnt"`
	Mgd              int                        `json:"mgd"`
	EventPoint       int                        `json:"event_point"`
	LiveDifficultyID int                        `json:"live_difficulty_id"`
	TimeStamp        int                        `json:"timeStamp"`
	PreciseScoreLog  liveschema.PreciseScoreLog `json:"precise_score_log"`
	ScoreCute        int                        `json:"score_cute"`
	EventID          any                        `json:"event_id"`
	ScoreCool        int                        `json:"score_cool"`
	Token            string                     `json:"token"`
}

func (req RewardReq) ToLiveRewardReq(liveDifficultyID int) liveschema.RewardReq {
	return liveschema.RewardReq{
		Module:           req.Module,
		Action:           req.Action,
		GoodCnt:          req.GoodCnt,
		MissCnt:          req.MissCnt,
		IsTraining:       req.IsTraining,
		GreatCnt:         req.GreatCnt,
		CommandNum:       req.CommandNum,
		LoveCnt:          req.LoveCnt,
		RemainHp:         req.RemainHp,
		MaxCombo:         req.MaxCombo,
		ScoreSmile:       req.ScoreSmile,
		PerfectCnt:       req.PerfectCnt,
		BadCnt:           req.BadCnt,
		Mgd:              req.Mgd,
		EventPoint:       req.EventPoint,
		LiveDifficultyID: liveDifficultyID,
		TimeStamp:        req.TimeStamp,
		PreciseScoreLog:  req.PreciseScoreLog,
		ScoreCute:        req.ScoreCute,
		EventID:          req.EventID,
		ScoreCool:        req.ScoreCool,
	}
}
