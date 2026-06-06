package liveschema

import (
	"encoding/json"
	"strconv"
)

type FlexibleLiveDifficultyID int

func (id *FlexibleLiveDifficultyID) UnmarshalJSON(data []byte) error {
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*id = FlexibleLiveDifficultyID(intValue)
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err != nil {
		return err
	}

	parsed, err := strconv.Atoi(stringValue)
	if err != nil {
		return err
	}
	*id = FlexibleLiveDifficultyID(parsed)
	return nil
}

type PlayScoreReq struct {
	Module           string                   `json:"module"`
	Action           string                   `json:"action"`
	TimeStamp        int64                    `json:"timeStamp"`
	Mgd              int                      `json:"mgd"`
	LiveDifficultyID FlexibleLiveDifficultyID `json:"live_difficulty_id"`
	CommandNum       string                   `json:"commandNum"`
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
	LiveSetting any      `json:"live_setting,omitempty"`
	CanReplay   bool     `json:"can_replay"`
	TriggerLog  any      `json:"trigger_log,omitempty"`
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
	LiveSetting any      `json:"live_setting,omitempty"`
	CanReplay   bool     `json:"can_replay"`
	TriggerLog  any      `json:"trigger_log,omitempty"`
}

type Skill struct {
	HasRecord   bool     `json:"has_record"`
	LiveInfo    LiveInfo `json:"live_info"`
	RandomSeed  any      `json:"random_seed"`
	MaxCombo    any      `json:"max_combo"`
	UpdateDate  any      `json:"update_date"`
	PreciseList any      `json:"precise_list"`
	DeckInfo    any      `json:"deck_info"`
	TapAdjust   any      `json:"tap_adjust"`
	LiveSetting any      `json:"live_setting,omitempty"`
	CanReplay   bool     `json:"can_replay"`
	TriggerLog  any      `json:"trigger_log,omitempty"`
}

type PreciseScoreData struct {
	On                Skill      `json:"on"`
	Off               Skill      `json:"off"`
	RankInfo          []RankInfo `json:"rank_info"`
	CanActivateEffect bool       `json:"can_activate_effect"`
	ServerTimestamp   int64      `json:"server_timestamp"`
}

type PreciseScoreResp struct {
	ResponseData any   `json:"response_data"`
	ReleaseInfo  []any `json:"release_info"`
	StatusCode   int   `json:"status_code"`
}
