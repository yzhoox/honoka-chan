package friendschema

type SearchReq struct {
	Module     string `json:"module"`
	Action     string `json:"action"`
	TimeStamp  int64  `json:"timeStamp"`
	Mgd        int    `json:"mgd"`
	InviteCode string `json:"invite_code"`
	CommandNum string `json:"commandNum"`
}

type SearchUserInfo struct {
	UserID               int    `json:"user_id"`
	Name                 string `json:"name"`
	Level                int    `json:"level"`
	CostMax              int    `json:"cost_max"`
	UnitMax              int    `json:"unit_max"`
	EnergyMax            int    `json:"energy_max"`
	FriendMax            int    `json:"friend_max"`
	UnitCnt              int    `json:"unit_cnt"`
	ElapsedTimeFromLogin string `json:"elapsed_time_from_login"`
	Comment              string `json:"comment"`
}

type SearchCostume struct {
	UnitID    int  `json:"unit_id"`
	IsRankMax bool `json:"is_rank_max"`
	IsSigned  bool `json:"is_signed"`
}

type SearchCenterUnitInfo struct {
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
	Costume                    SearchCostume `json:"costume"`
}

type SearchData struct {
	UserInfo        SearchUserInfo       `json:"user_info"`
	CenterUnitInfo  SearchCenterUnitInfo `json:"center_unit_info"`
	SettingAwardID  int                  `json:"setting_award_id"`
	IsAlliance      bool                 `json:"is_alliance"`
	FriendStatus    int                  `json:"friend_status"`
	ServerTimestamp int64                `json:"server_timestamp"`
}

type SearchResp struct {
	ResponseData SearchData `json:"response_data"`
	ReleaseInfo  []any      `json:"release_info"`
	StatusCode   int        `json:"status_code"`
}
