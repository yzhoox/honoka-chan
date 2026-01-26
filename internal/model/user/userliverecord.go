package usermodel

type UserLiveRecord struct {
	ID               int    `xorm:"id pk autoincr"`
	LiveDifficultyID int    `xorm:"live_difficulty_id"`
	DeckID           int    `xorm:"deck_id"`
	TotalScore       int    `xorm:"total_score"` //
	PerfectCnt       int    `xorm:"perfect_cnt"` //
	GreatCnt         int    `xorm:"great_cnt"`   //
	GoodCnt          int    `xorm:"good_cnt"`    //
	BadCnt           int    `xorm:"bad_cnt"`     //
	MissCnt          int    `xorm:"miss_cnt"`    //
	MaxCombo         int    `xorm:"max_combo"`
	IsSkillOn        bool   `xorm:"is_skill_on"`
	LiveInfoJSON     string `xorm:"live_info_json"`
	PreciseListJSON  string `xorm:"precise_list_json"`
	DeckInfoJSON     string `xorm:"deck_info_json"`
	TapAdjust        int    `xorm:"tap_adjust"`
	LiveSettingJSON  string `xorm:"live_setting_json"`
	CanReplay        bool   `xorm:"can_replay"`
	TriggerLogJSON   string `xorm:"trigger_log_json"`
	UserID           int    `xorm:"user_id"`
	UpdateDate       string `xorm:"insert_date"`
}

// TODO: 转移到 schema 下
type NotesList struct {
	TimingSec      float64 `json:"timing_sec"`
	NotesAttribute int     `json:"notes_attribute"`
	NotesLevel     int     `json:"notes_level"`
	Effect         int     `json:"effect"`
	EffectValue    int     `json:"effect_value"`
	Position       int     `json:"position"`
}

type LiveInfo struct {
	LiveDifficultyID int         `json:"live_difficulty_id"`
	IsRandom         bool        `json:"is_random"`
	AcFlag           int         `json:"ac_flag"`
	SwingFlag        int         `json:"swing_flag"`
	NotesList        []NotesList `json:"notes_list"`
}

// TODO: 转移到 schema 下
type PreciseList struct {
	Count      int     `json:"count"`
	Effect     int     `json:"effect"`
	Accuracy   int     `json:"accuracy"`
	Tap        float64 `json:"tap"`
	IsSame     bool    `json:"is_same"`
	Release    any     `json:"release"`
	Position   int     `json:"position"`
	NoteNumber int     `json:"note_number"`
	FirstTouch any     `json:"first_touch"`
	Tp         bool    `json:"tp"`
	Tpf        any     `json:"tpf"`
}

// TODO: 转移到 schema 下
type TotalStatus struct {
	Hp    int `json:"hp"`
	Smile int `json:"smile"`
	Cute  int `json:"cute"`
	Cool  int `json:"cool"`
}

type CenterBonus struct {
	Hp    int `json:"hp"`
	Smile int `json:"smile"`
	Cute  int `json:"cute"`
	Cool  int `json:"cool"`
}

type SiBonus struct {
	Hp    int `json:"hp"`
	Smile int `json:"smile"`
	Cute  int `json:"cute"`
	Cool  int `json:"cool"`
}

type AccessoryInfo struct {
	AccessoryOwningUserID int  `json:"accessory_owning_user_id"`
	AccessoryID           int  `json:"accessory_id"`
	Exp                   int  `json:"exp"`
	NextExp               int  `json:"next_exp"`
	Level                 int  `json:"level"`
	MaxLevel              int  `json:"max_level"`
	RankUpCount           int  `json:"rank_up_count"`
	FavoriteFlag          bool `json:"favorite_flag"`
}

type UnitList struct {
	UnitID                     int            `json:"unit_id"`
	Position                   int            `json:"position"`
	Level                      int            `json:"level"`
	LevelLimitID               int            `json:"level_limit_id"`
	DisplayRank                int            `json:"display_rank"`
	Love                       int            `json:"love"`
	UnitSkillLevel             int            `json:"unit_skill_level"`
	IsRankMax                  bool           `json:"is_rank_max"`
	IsLoveMax                  bool           `json:"is_love_max"`
	IsLevelMax                 bool           `json:"is_level_max"`
	IsSigned                   bool           `json:"is_signed"`
	Rank                       int            `json:"rank"`
	MaxHp                      int            `json:"max_hp"`
	UnitRemovableSkillCapacity int            `json:"unit_removable_skill_capacity"`
	RemovableSkillIds          []int          `json:"removable_skill_ids"`
	TotalStatus                TotalStatus    `json:"total_status"`
	SiBonus                    SiBonus        `json:"si_bonus"`
	AccessoryInfo              *AccessoryInfo `json:"accessory_info,omitempty"`
}

type DeckInfo struct {
	LiveDifficultyID int         `json:"live_difficulty_id"`
	TotalStatus      TotalStatus `json:"total_status"`
	CenterBonus      CenterBonus `json:"center_bonus"`
	SiBonus          SiBonus     `json:"si_bonus"`
	UnitList         []UnitList  `json:"unit_list"`
}

// TODO: 转移到 schema 下
type Icon struct {
	NormalID int `json:"normal_id"`
	JustID   int `json:"just_id"`
	SlideID  int `json:"slide_id"`
}

type LiveSetting struct {
	NotesSpeed      float64 `json:"notes_speed"`
	CutinBrightness int     `json:"cutin_brightness"`
	CutinType       int     `json:"cutin_type"`
	EffectFlag      bool    `json:"effect_flag"`
	StringSize      int     `json:"string_size"`
	SeID            int     `json:"se_id"`
	Icon            Icon    `json:"icon"`
	BackgroundID    int     `json:"background_id"`
	RandomValue     int     `json:"random_value"`
}

// TODO: 转移到 schema 下
type TriggerLog struct {
	Position       int `json:"position"`
	ActivationRate int `json:"activation_rate"`
}
