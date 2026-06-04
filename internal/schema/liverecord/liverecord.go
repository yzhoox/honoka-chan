package liverecordschema

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

type TriggerLog struct {
	Position       int `json:"position"`
	ActivationRate int `json:"activation_rate"`
}
