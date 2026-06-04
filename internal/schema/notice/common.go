package noticeschema

type GreetingNoticeReq struct {
	Module     string `json:"module"`
	Action     string `json:"action"`
	TimeStamp  int64  `json:"timeStamp"`
	Mgd        int    `json:"mgd"`
	NextID     int    `json:"next_id"`
	IsPrevious bool   `json:"is_previous"`
	CommandNum string `json:"commandNum"`
}

type GreetingUserData struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Level  int    `json:"level"`
}

type GreetingAccessoryInfo struct {
	AccessoryOwningUserID int  `json:"accessory_owning_user_id"`
	AccessoryID           int  `json:"accessory_id"`
	Exp                   int  `json:"exp"`
	NextExp               int  `json:"next_exp"`
	Level                 int  `json:"level"`
	MaxLevel              int  `json:"max_level"`
	RankUpCount           int  `json:"rank_up_count"`
	FavoriteFlag          bool `json:"favorite_flag"`
}

type GreetingCenterUnitInfo struct {
	UnitOwningUserID           int64                  `json:"unit_owning_user_id"`
	UnitID                     int                    `json:"unit_id"`
	Exp                        int                    `json:"exp"`
	NextExp                    int                    `json:"next_exp"`
	Level                      int                    `json:"level"`
	LevelLimitID               int                    `json:"level_limit_id"`
	MaxLevel                   int                    `json:"max_level"`
	Rank                       int                    `json:"rank"`
	MaxRank                    int                    `json:"max_rank"`
	Love                       int                    `json:"love"`
	MaxLove                    int                    `json:"max_love"`
	UnitSkillLevel             int                    `json:"unit_skill_level"`
	MaxHp                      int                    `json:"max_hp"`
	FavoriteFlag               bool                   `json:"favorite_flag"`
	DisplayRank                int                    `json:"display_rank"`
	UnitSkillExp               int                    `json:"unit_skill_exp"`
	UnitRemovableSkillCapacity int                    `json:"unit_removable_skill_capacity"`
	Attribute                  int                    `json:"attribute"`
	Smile                      int                    `json:"smile"`
	Cute                       int                    `json:"cute"`
	Cool                       int                    `json:"cool"`
	IsLoveMax                  bool                   `json:"is_love_max"`
	IsLevelMax                 bool                   `json:"is_level_max"`
	IsRankMax                  bool                   `json:"is_rank_max"`
	IsSigned                   bool                   `json:"is_signed"`
	IsSkillLevelMax            bool                   `json:"is_skill_level_max"`
	SettingAwardID             int                    `json:"setting_award_id"`
	RemovableSkillIds          []int                  `json:"removable_skill_ids"`
	AccessoryInfo              *GreetingAccessoryInfo `json:"accessory_info,omitempty"`
}

type GreetingPeer struct {
	UserData       GreetingUserData       `json:"user_data"`
	CenterUnitInfo GreetingCenterUnitInfo `json:"center_unit_info"`
	SettingAwardID int                    `json:"setting_award_id"`
}
