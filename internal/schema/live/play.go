package live

type PlayReq struct {
	Module           string `json:"module"`
	PartyUserID      int64  `json:"party_user_id"`
	Action           string `json:"action"`
	Mgd              int    `json:"mgd"`
	IsTraining       bool   `json:"is_training"`
	UnitDeckID       int    `json:"unit_deck_id"`
	LiveDifficultyID string `json:"live_difficulty_id"`
	TimeStamp        int    `json:"timeStamp"`
	LpFactor         int    `json:"lp_factor"`
	CommandNum       string `json:"commandNum"`
}

type RankInfo struct {
	Rank    int `json:"rank"`
	RankMin int `json:"rank_min"`
	RankMax int `json:"rank_max"`
}

type NotesList struct {
	TimingSec      float64 `json:"timing_sec"`
	NotesAttribute int     `json:"notes_attribute"`
	NotesLevel     int     `json:"notes_level"`
	Effect         int     `json:"effect"`
	EffectValue    float64 `json:"effect_value"`
	Position       int     `json:"position"`
}

type LiveInfo struct {
	LiveDifficultyID int         `json:"live_difficulty_id"`
	IsRandom         bool        `json:"is_random"`
	AcFlag           int         `json:"ac_flag"`
	SwingFlag        int         `json:"swing_flag"`
	NotesList        []NotesList `json:"notes_list"`
}

type Costume struct {
	UnitID    int  `json:"unit_id"`
	IsRankMax bool `json:"is_rank_max"`
	IsSigned  bool `json:"is_signed"`
}

type UnitList struct {
	Smile int `json:"smile"`
	Cute  int `json:"cute"`
	Cool  int `json:"cool"`
	// Costume Costume `json:"costume"`
}

type DeckInfo struct {
	UnitDeckID       int        `json:"unit_deck_id"`
	TotalSmile       int        `json:"total_smile"`
	TotalCute        int        `json:"total_cute"`
	TotalCool        int        `json:"total_cool"`
	TotalHp          int        `json:"total_hp"`
	PreparedHpDamage int        `json:"prepared_hp_damage"`
	UnitList         []UnitList `json:"unit_list"`
}

type PlayLiveList struct {
	LiveInfo LiveInfo `json:"live_info"`
	DeckInfo DeckInfo `json:"deck_info"`
}

type PlayData struct {
	RankInfo            []RankInfo     `json:"rank_info"`
	EnergyFullTime      string         `json:"energy_full_time"`
	OverMaxEnergy       int            `json:"over_max_energy"`
	AvailableLiveResume bool           `json:"available_live_resume"`
	LiveList            []PlayLiveList `json:"live_list"`
	IsMarathonEvent     bool           `json:"is_marathon_event"`
	MarathonEventID     any            `json:"marathon_event_id"`
	NoSkill             bool           `json:"no_skill"`
	CanActivateEffect   bool           `json:"can_activate_effect"`
	ServerTimestamp     int64          `json:"server_timestamp"`
}

type PlayResp struct {
	ResponseData PlayData `json:"response_data"`
	ReleaseInfo  []any    `json:"release_info"`
	StatusCode   int      `json:"status_code"`
}
