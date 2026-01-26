package liveschema

import userschema "honoka-chan/internal/schema/user"

type CenterUnitInfo struct {
	UnitOwningUserID           int   `json:"unit_owning_user_id"`
	UnitID                     int   `json:"unit_id"`
	Exp                        int   `json:"exp"`
	NextExp                    int   `json:"next_exp"`
	Level                      int   `json:"level"`
	LevelLimitID               int   `json:"level_limit_id"`
	MaxLevel                   int   `json:"max_level"`
	Rank                       int   `json:"rank"`
	MaxRank                    int   `json:"max_rank"`
	Love                       int   `json:"love"`
	MaxLove                    int   `json:"max_love"`
	UnitSkillLevel             int   `json:"unit_skill_level"`
	MaxHp                      int   `json:"max_hp"`
	FavoriteFlag               bool  `json:"favorite_flag"`
	DisplayRank                int   `json:"display_rank"`
	UnitSkillExp               int   `json:"unit_skill_exp"`
	UnitRemovableSkillCapacity int   `json:"unit_removable_skill_capacity"`
	Attribute                  int   `json:"attribute"`
	Smile                      int   `json:"smile"`
	Cute                       int   `json:"cute"`
	Cool                       int   `json:"cool"`
	IsLoveMax                  bool  `json:"is_love_max"`
	IsLevelMax                 bool  `json:"is_level_max"`
	IsRankMax                  bool  `json:"is_rank_max"`
	IsSigned                   bool  `json:"is_signed"`
	IsSkillLevelMax            bool  `json:"is_skill_level_max"`
	SettingAwardID             int   `json:"setting_award_id"`
	RemovableSkillIds          []int `json:"removable_skill_ids"`
}

type PartyList struct {
	UserInfo             userschema.UserInfo `json:"user_info"`
	CenterUnitInfo       CenterUnitInfo      `json:"center_unit_info"`
	SettingAwardID       int                 `json:"setting_award_id"`
	AvailableSocialPoint int                 `json:"available_social_point"`
	FriendStatus         int                 `json:"friend_statuses"`
}

type PartyListData struct {
	PartyList         []PartyList `json:"party_list"`
	TrainingEnergy    int         `json:"training_energy"`
	TrainingEnergyMax int         `json:"training_energy_max"`
	ServerTimestamp   int64       `json:"server_timestamp"`
}

type PartyListResp struct {
	ResponseData PartyListData `json:"response_data"`
	ReleaseInfo  []any         `json:"release_info"`
	StatusCode   int           `json:"status_code"`
}
