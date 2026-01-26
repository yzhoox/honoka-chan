package unitmodel

type UnitM struct {
	UnitId                        int     `xorm:"unit_id pk"`
	UnitNumber                    int     `xorm:"unit_number"`
	UnitTypeId                    int     `xorm:"unit_type_id"`
	AlbumSeriesId                 *int    `xorm:"album_series_id"`
	Eponym                        *string `xorm:"eponym"`
	EponymEn                      string  `xorm:"eponym_en"`
	Name                          string  `xorm:"name"`
	NameEn                        *string `xorm:"name_en"`
	NormalCardId                  int     `xorm:"normal_card_id"`
	RankMaxCardId                 int     `xorm:"rank_max_card_id"`
	NormalIconAsset               string  `xorm:"normal_icon_asset"`
	NormalIconAssetEn             *string `xorm:"normal_icon_asset_en"`
	RankMaxIconAsset              string  `xorm:"rank_max_icon_asset"`
	RankMaxIconAssetEn            *string `xorm:"rank_max_icon_asset_en"`
	NormalUnitNaviAssetId         int     `xorm:"normal_unit_navi_asset_id"`
	RankMaxUnitNaviAssetId        int     `xorm:"rank_max_unit_navi_asset_id"`
	Rarity                        int     `xorm:"rarity"`
	AttributeId                   int     `xorm:"attribute_id"`
	DefaultUnitSkillId            *int    `xorm:"default_unit_skill_id"`
	SkillAssetVoiceId             *int    `xorm:"skill_asset_voice_id"`
	DefaultLeaderSkillId          *int    `xorm:"default_leader_skill_id"`
	DefaultRemovableSkillCapacity int     `xorm:"default_removable_skill_capacity"`
	MaxRemovableSkillCapacity     int     `xorm:"max_removable_skill_capacity"`
	DisableRankUp                 int     `xorm:"disable_rank_up"`
	RankMin                       int     `xorm:"rank_min"`
	RankMax                       int     `xorm:"rank_max"`
	UnitLevelUpPatternId          int     `xorm:"unit_level_up_pattern_id"`
	HpMax                         int     `xorm:"hp_max"`
	SmileMax                      int     `xorm:"smile_max"`
	PureMax                       int     `xorm:"pure_max"`
	CoolMax                       int     `xorm:"cool_max"`
	ReinforceItemRankUpCost       *int    `xorm:"reinforce_item_rank_up_cost"`
	SubUnitTypeId                 *int    `xorm:"sub_unit_type_id"`
	ReleaseTag                    *string `xorm:"release_tag"`
}
