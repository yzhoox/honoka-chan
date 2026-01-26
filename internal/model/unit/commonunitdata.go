package unitmodel

type CommonUnitData struct {
	ID                          int    `xorm:"id pk autoincr" json:"-"`
	UnitNumber                  int    `xorm:"unit_number" json:"-"`
	UnitID                      int    `xorm:"unit_id" json:"unit_id"`
	UnitTypeID                  int    `xorm:"unit_type_id" json:"unit_type_id"`
	Name                        string `xorm:"name" json:"name"`
	Eponym                      string `xorm:"eponym" json:"eponym"`
	Rarity                      int    `xorm:"rarity" json:"rarity"`
	Attribute                   int    `xorm:"attribute" json:"attribute"`
	Smile                       int    `xorm:"smile" json:"smile"`
	Cute                        int    `xorm:"cute" json:"cute"`
	Cool                        int    `xorm:"cool" json:"cool"`
	Exp                         int    `xorm:"exp" json:"exp"`
	Level                       int    `xorm:"level" json:"level"`
	MaxLevel                    int    `xorm:"max_level" json:"max_level"`
	LevelLimitID                int    `xorm:"level_limit_id" json:"level_limit_id"`
	Rank                        int    `xorm:"rank" json:"rank"`
	MaxRank                     int    `xorm:"max_rank" json:"max_rank"`
	Love                        int    `xorm:"love" json:"love"`
	MaxLove                     int    `xorm:"max_love" json:"max_love"`
	UnitSkillExp                int    `xorm:"unit_skill_exp" json:"unit_skill_exp"`
	UnitSkillLevel              int    `xorm:"unit_skill_level" json:"unit_skill_level"`
	MaxHp                       int    `xorm:"max_hp" json:"max_hp"`
	UnitRemovableSkillCapacity  int    `xorm:"unit_removable_skill_capacity" json:"unit_removable_skill_capacity"`
	IsRankMax                   bool   `xorm:"is_rank_max" json:"is_rank_max"`
	IsLoveMax                   bool   `xorm:"is_love_max" json:"is_love_max"`
	IsLevelMax                  bool   `xorm:"is_level_max" json:"is_level_max"`
	IsSigned                    bool   `xorm:"is_signed" json:"is_signed"`
	IsSkillLevelMax             bool   `xorm:"is_skill_level_max" json:"is_skill_level_max"`
	IsRemovableSkillCapacityMax bool   `xorm:"is_removable_skill_capacity_max" json:"is_removable_skill_capacity_max"`
	InsertDate                  int64  `xorm:"insert_date" json:"insert_date"`
}

type UnitDataMap struct {
	UnitOwningUserID int  `xorm:"unit_owning_user_id pk autoincr"`
	FavoriteFlag     bool `xorm:"favorite_flag"`
	DisplayRank      int  `xorm:"display_rank"`

	ID                          int    `xorm:"id pk autoincr" json:"-"`
	UnitID                      int    `xorm:"unit_id" json:"unit_id"`
	UnitTypeID                  int    `xorm:"unit_type_id" json:"unit_type_id"`
	Name                        string `xorm:"name" json:"name"`
	Eponym                      string `xorm:"eponym" json:"eponym"`
	Rarity                      int    `xorm:"rarity" json:"rarity"`
	Attribute                   int    `xorm:"attribute" json:"attribute"`
	Smile                       int    `xorm:"smile" json:"smile"`
	Cute                        int    `xorm:"cute" json:"cute"`
	Cool                        int    `xorm:"cool" json:"cool"`
	Exp                         int    `xorm:"exp" json:"exp"`
	Level                       int    `xorm:"level" json:"level"`
	MaxLevel                    int    `xorm:"max_level" json:"max_level"`
	LevelLimitID                int    `xorm:"level_limit_id" json:"level_limit_id"`
	Rank                        int    `xorm:"rank" json:"rank"`
	MaxRank                     int    `xorm:"max_rank" json:"max_rank"`
	Love                        int    `xorm:"love" json:"love"`
	MaxLove                     int    `xorm:"max_love" json:"max_love"`
	UnitSkillExp                int    `xorm:"unit_skill_exp" json:"unit_skill_exp"`
	UnitSkillLevel              int    `xorm:"unit_skill_level" json:"unit_skill_level"`
	MaxHp                       int    `xorm:"max_hp" json:"max_hp"`
	UnitRemovableSkillCapacity  int    `xorm:"unit_removable_skill_capacity" json:"unit_removable_skill_capacity"`
	IsRankMax                   bool   `xorm:"is_rank_max" json:"is_rank_max"`
	IsLoveMax                   bool   `xorm:"is_love_max" json:"is_love_max"`
	IsLevelMax                  bool   `xorm:"is_level_max" json:"is_level_max"`
	IsSigned                    bool   `xorm:"is_signed" json:"is_signed"`
	IsSkillLevelMax             bool   `xorm:"is_skill_level_max" json:"is_skill_level_max"`
	IsRemovableSkillCapacityMax bool   `xorm:"is_removable_skill_capacity_max" json:"is_removable_skill_capacity_max"`
	InsertDate                  int64  `xorm:"insert_date" json:"insert_date"`
}
