package rliveschema

import liveschema "honoka-chan/internal/schema/live"

type PlayReq struct {
	Module      string `json:"module"`
	PartyUserID int64  `json:"party_user_id"`
	Action      string `json:"action"`
	Mgd         int    `json:"mgd"`
	IsTraining  bool   `json:"is_training"`
	UnitDeckID  int    `json:"unit_deck_id"`
	Token       string `json:"token"`
	TimeStamp   int    `json:"timeStamp"`
	LpFactor    int    `json:"lp_factor"`
	CommandNum  string `json:"commandNum"`
}

func (req PlayReq) ToLivePlayReq(liveDifficultyID int) liveschema.PlayReq {
	return liveschema.PlayReq{
		Module:           req.Module,
		PartyUserID:      req.PartyUserID,
		Action:           req.Action,
		Mgd:              req.Mgd,
		IsTraining:       req.IsTraining,
		UnitDeckID:       req.UnitDeckID,
		LiveDifficultyID: itoa(liveDifficultyID),
		TimeStamp:        req.TimeStamp,
		LpFactor:         req.LpFactor,
		CommandNum:       req.CommandNum,
	}
}
