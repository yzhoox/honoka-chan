package rliveschema

import liveschema "honoka-chan/internal/schema/live"

type LotReq struct {
	Module         string `json:"module"`
	Action         string `json:"action"`
	MemberCategory int    `json:"member_category"`
	Mgd            int    `json:"mgd"`
	Difficulty     int    `json:"difficulty"`
	Attribute      int    `json:"attribute"`
	TimeStamp      int64  `json:"timeStamp"`
	CommandNum     string `json:"commandNum"`
}

type LiveInfo struct {
	LiveDifficultyID int  `json:"live_difficulty_id"`
	IsRandom         bool `json:"is_random"`
	AcFlag           int  `json:"ac_flag"`
	SwingFlag        int  `json:"swing_flag"`
}

type LotData struct {
	LiveInfo          LiveInfo               `json:"live_info"`
	HasSlideNotes     int                    `json:"has_slide_notes"`
	PartyList         []liveschema.PartyList `json:"party_list"`
	TrainingEnergy    int                    `json:"training_energy"`
	TrainingEnergyMax int                    `json:"training_energy_max"`
	Token             string                 `json:"token"`
}

type LotResp struct {
	ResponseData LotData `json:"response_data"`
	ReleaseInfo  []any   `json:"release_info"`
	StatusCode   int     `json:"status_code"`
}
