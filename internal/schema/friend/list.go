package friendschema

type ListReq struct {
	Module     string `json:"module"`
	Type       int    `json:"type"`
	Action     string `json:"action"`
	TimeStamp  int64  `json:"timeStamp"`
	Mgd        int    `json:"mgd"`
	Sort       int    `json:"sort"`
	Page       int    `json:"page"`
	CommandNum string `json:"commandNum"`
}

type UserData struct {
	UserID                 int    `json:"user_id"`
	Name                   string `json:"name"`
	Level                  int    `json:"level"`
	ElapsedTimeFromLogin   string `json:"elapsed_time_from_login"`
	ElapsedTimeFromApplied string `json:"elapsed_time_from_applied"`
	Comment                string `json:"comment"`
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

type CenterUnitInfo struct {
	UnitOwningUserID           int64         `json:"unit_owning_user_id"`
	UnitID                     int           `json:"unit_id"`
	Exp                        int           `json:"exp"`
	NextExp                    int           `json:"next_exp"`
	Level                      int           `json:"level"`
	LevelLimitID               int           `json:"level_limit_id"`
	MaxLevel                   int           `json:"max_level"`
	Rank                       int           `json:"rank"`
	MaxRank                    int           `json:"max_rank"`
	Love                       int           `json:"love"`
	MaxLove                    int           `json:"max_love"`
	UnitSkillLevel             int           `json:"unit_skill_level"`
	MaxHp                      int           `json:"max_hp"`
	FavoriteFlag               bool          `json:"favorite_flag"`
	DisplayRank                int           `json:"display_rank"`
	UnitSkillExp               int           `json:"unit_skill_exp"`
	UnitRemovableSkillCapacity int           `json:"unit_removable_skill_capacity"`
	Attribute                  int           `json:"attribute"`
	Smile                      int           `json:"smile"`
	Cute                       int           `json:"cute"`
	Cool                       int           `json:"cool"`
	IsLoveMax                  bool          `json:"is_love_max"`
	IsLevelMax                 bool          `json:"is_level_max"`
	IsRankMax                  bool          `json:"is_rank_max"`
	IsSigned                   bool          `json:"is_signed"`
	IsSkillLevelMax            bool          `json:"is_skill_level_max"`
	SettingAwardID             int           `json:"setting_award_id"`
	RemovableSkillIds          []int         `json:"removable_skill_ids"`
	AccessoryInfo              AccessoryInfo `json:"accessory_info"`
}

type FriendList struct {
	UserData       UserData       `json:"user_data"`
	CenterUnitInfo CenterUnitInfo `json:"center_unit_info"`
	SettingAwardID int            `json:"setting_award_id"`
	IsNew          bool           `json:"is_new"`
}

type ListData struct {
	ItemCount       int          `json:"item_count"`
	FriendList      []FriendList `json:"friend_list"`
	NewFriendList   []any        `json:"new_friend_list"`
	ServerTimestamp int64        `json:"server_timestamp"`
}

type ListResp struct {
	ResponseData ListData `json:"response_data"`
	ReleaseInfo  []any    `json:"release_info"`
	StatusCode   int      `json:"status_code"`
}
