package live

type PlayScoreReq struct {
	Module           string `json:"module"`
	Action           string `json:"action"`
	TimeStamp        int64  `json:"timeStamp"`
	Mgd              int    `json:"mgd"`
	LiveDifficultyID string `json:"live_difficulty_id"`
	CommandNum       string `json:"commandNum"`
}

type On struct {
	HasRecord   bool     `json:"has_record"`
	LiveInfo    LiveInfo `json:"live_info"`
	RandomSeed  any      `json:"random_seed"`
	MaxCombo    any      `json:"max_combo"`
	UpdateDate  any      `json:"update_date"`
	PreciseList any      `json:"precise_list"`
	DeckInfo    any      `json:"deck_info"`
	TapAdjust   any      `json:"tap_adjust"`
	CanReplay   bool     `json:"can_replay"`
}

type Off struct {
	HasRecord   bool     `json:"has_record"`
	LiveInfo    LiveInfo `json:"live_info"`
	RandomSeed  any      `json:"random_seed"`
	MaxCombo    any      `json:"max_combo"`
	UpdateDate  any      `json:"update_date"`
	PreciseList any      `json:"precise_list"`
	DeckInfo    any      `json:"deck_info"`
	TapAdjust   any      `json:"tap_adjust"`
	CanReplay   bool     `json:"can_replay"`
}

type PreciseScoreData struct {
	On                On         `json:"on"`
	Off               Off        `json:"off"`
	RankInfo          []RankInfo `json:"rank_info"`
	CanActivateEffect bool       `json:"can_activate_effect"`
	ServerTimestamp   int64      `json:"server_timestamp"`
}

type PreciseScoreResp struct {
	ResponseData PreciseScoreData `json:"response_data"`
	ReleaseInfo  []any            `json:"release_info"`
	StatusCode   int              `json:"status_code"`
}
