package unitmodel

type CommonUnitData struct {
	ID                          int    `xorm:"id pk autoincr"`
	UnitNumber                  int    `xorm:"unit_number"`
	UnitID                      int    `xorm:"unit_id"`
	UnitTypeID                  int    `xorm:"unit_type_id"`
	Name                        string `xorm:"name"`
	Eponym                      string `xorm:"eponym"`
	Rarity                      int    `xorm:"rarity"`
	Attribute                   int    `xorm:"attribute"`
	Smile                       int    `xorm:"smile"`
	Cute                        int    `xorm:"cute"`
	Cool                        int    `xorm:"cool"`
	Exp                         int    `xorm:"exp"`
	Level                       int    `xorm:"level"`
	MaxLevel                    int    `xorm:"max_level"`
	LevelLimitID                int    `xorm:"level_limit_id"`
	Rank                        int    `xorm:"rank"`
	MaxRank                     int    `xorm:"max_rank"`
	Love                        int    `xorm:"love"`
	MaxLove                     int    `xorm:"max_love"`
	UnitSkillExp                int    `xorm:"unit_skill_exp"`
	UnitSkillLevel              int    `xorm:"unit_skill_level"`
	MaxHp                       int    `xorm:"max_hp"`
	UnitRemovableSkillCapacity  int    `xorm:"unit_removable_skill_capacity"`
	IsRankMax                   bool   `xorm:"is_rank_max"`
	IsLoveMax                   bool   `xorm:"is_love_max"`
	IsLevelMax                  bool   `xorm:"is_level_max"`
	IsSigned                    bool   `xorm:"is_signed"`
	IsSkillLevelMax             bool   `xorm:"is_skill_level_max"`
	IsRemovableSkillCapacityMax bool   `xorm:"is_removable_skill_capacity_max"`
	InsertDate                  int64  `xorm:"insert_date"`
}

type UnitDataMap struct {
	ID                          int    `xorm:"id pk autoincr"`
	UnitOwningUserID            int    `xorm:"unit_owning_user_id pk autoincr"`
	UnitID                      int    `xorm:"unit_id"`
	UnitTypeID                  int    `xorm:"unit_type_id"`
	Name                        string `xorm:"name"`
	Eponym                      string `xorm:"eponym"`
	Rarity                      int    `xorm:"rarity"`
	Attribute                   int    `xorm:"attribute"`
	Smile                       int    `xorm:"smile"`
	Cute                        int    `xorm:"cute"`
	Cool                        int    `xorm:"cool"`
	Exp                         int    `xorm:"exp"`
	Level                       int    `xorm:"level"`
	MaxLevel                    int    `xorm:"max_level"`
	LevelLimitID                int    `xorm:"level_limit_id"`
	Rank                        int    `xorm:"rank"`
	MaxRank                     int    `xorm:"max_rank"`
	DisplayRank                 int    `xorm:"display_rank"`
	Love                        int    `xorm:"love"`
	MaxLove                     int    `xorm:"max_love"`
	UnitSkillExp                int    `xorm:"unit_skill_exp"`
	UnitSkillLevel              int    `xorm:"unit_skill_level"`
	MaxHp                       int    `xorm:"max_hp"`
	UnitRemovableSkillCapacity  int    `xorm:"unit_removable_skill_capacity"`
	IsRankMax                   bool   `xorm:"is_rank_max"`
	IsLoveMax                   bool   `xorm:"is_love_max"`
	IsLevelMax                  bool   `xorm:"is_level_max"`
	IsSigned                    bool   `xorm:"is_signed"`
	IsSkillLevelMax             bool   `xorm:"is_skill_level_max"`
	IsRemovableSkillCapacityMax bool   `xorm:"is_removable_skill_capacity_max"`
	FavoriteFlag                bool   `xorm:"favorite_flag"`
	InsertDate                  int64  `xorm:"insert_date"`
}
