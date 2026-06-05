package secretboxschema

import secretboxapischema "honoka-chan/internal/schema/api/secretbox"

type PonReq struct {
	Module      string `json:"module"`
	Mgd         int    `json:"mgd"`
	Action      string `json:"action"`
	TimeStamp   int64  `json:"timeStamp"`
	ID          int    `json:"id"`
	CommandNum  string `json:"commandNum"`
	SecretBoxID int    `json:"secret_box_id"`
	UnitTypeIds []any  `json:"unit_type_ids"`
}

type SecretBoxUnitItem struct {
	UnitOwningUserID            int    `json:"unit_owning_user_id"`
	UnitOwningIDs               []int  `json:"unit_owning_ids"`
	UnitID                      int    `json:"unit_id"`
	Exp                         int    `json:"exp"`
	NextExp                     int    `json:"next_exp"`
	Level                       int    `json:"level"`
	MaxLevel                    int    `json:"max_level"`
	LevelLimitID                int    `json:"level_limit_id"`
	Rank                        int    `json:"rank"`
	MaxRank                     int    `json:"max_rank"`
	Love                        int    `json:"love"`
	MaxLove                     int    `json:"max_love"`
	UnitSkillLevel              int    `json:"unit_skill_level"`
	SkillLevel                  int    `json:"skill_level"`
	UnitSkillExp                int    `json:"unit_skill_exp"`
	MaxHp                       int    `json:"max_hp"`
	FavoriteFlag                bool   `json:"favorite_flag"`
	DisplayRank                 int    `json:"display_rank"`
	UnitRemovableSkillCapacity  int    `json:"unit_removable_skill_capacity"`
	MaxRemovableSkillCapacity   int    `json:"max_removable_skill_capacity"`
	Attribute                   int    `json:"attribute"`
	Smile                       int    `json:"smile"`
	Cute                        int    `json:"cute"`
	Cool                        int    `json:"cool"`
	IsRankMax                   bool   `json:"is_rank_max"`
	IsLoveMax                   bool   `json:"is_love_max"`
	IsLevelMax                  bool   `json:"is_level_max"`
	IsSigned                    bool   `json:"is_signed"`
	IsSkillLevelMax             bool   `json:"is_skill_level_max"`
	IsRemovableSkillCapacityMax bool   `json:"is_removable_skill_capacity_max"`
	IsSupportMember             bool   `json:"is_support_member"`
	NewUnitFlag                 bool   `json:"new_unit_flag"`
	UnitRarityID                int    `json:"unit_rarity_id"`
	IsHit                       bool   `json:"is_hit"`
	RemovableSkillIDs           []int  `json:"removable_skill_ids"`
	InsertDate                  string `json:"insert_date"`
	TotalSmile                  int    `json:"total_smile"`
	TotalCute                   int    `json:"total_cute"`
	TotalCool                   int    `json:"total_cool"`
	TotalHp                     int    `json:"total_hp"`
}

type SecretBoxRewardItem struct {
	ItemID         int    `json:"item_id"`
	AddType        int    `json:"add_type"`
	Amount         int    `json:"amount"`
	ItemCategoryID int    `json:"item_category_id"`
	RewardBoxFlag  bool   `json:"reward_box_flag"`
	InsertDate     string `json:"insert_date"`
}

type SecretBoxItems struct {
	Unit []SecretBoxUnitItem   `json:"unit"`
	Item []SecretBoxRewardItem `json:"item"`
}

type PonData struct {
	IsUnitMax                bool                                 `json:"is_unit_max"`
	ItemList                 []secretboxapischema.ItemAmount      `json:"item_list"`
	GaugeInfo                secretboxapischema.GaugeInfo         `json:"gauge_info"`
	ButtonList               []secretboxapischema.SecretBoxButton `json:"button_list"`
	SecretBoxInfo            secretboxapischema.SecretBoxInfo     `json:"secret_box_info"`
	SecretBoxItems           SecretBoxItems                       `json:"secret_box_items"`
	BeforeUserInfo           any                                  `json:"before_user_info"`
	AfterUserInfo            any                                  `json:"after_user_info"`
	FreeMuseGachaFlag        bool                                 `json:"free_muse_gacha_flag"`
	FreeAqoursGachaFlag      bool                                 `json:"free_aqours_gacha_flag"`
	LowestRarity             int                                  `json:"lowest_rarity"`
	PromotionPerformanceRate int                                  `json:"promotion_performance_rate"`
	SecretBoxParcelType      int                                  `json:"secret_box_parcel_type"`
	LimitBonusInfo           []any                                `json:"limit_bonus_info"`
	LimitBonusRewards        []any                                `json:"limit_bonus_rewards"`
	UnitSupportList          []any                                `json:"unit_support_list"`
}

type PonResp struct {
	ResponseData PonData `json:"response_data"`
	ReleaseInfo  []any   `json:"release_info"`
	StatusCode   int     `json:"status_code"`
}
