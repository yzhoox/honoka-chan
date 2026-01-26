package usermodel

type UserUnit struct {
	ID                          int    `xorm:"id pk autoincr"`
	UnitOwningUserID            int    `xorm:"unit_owning_user_id"`
	UserID                      int    `xorm:"user_id"`
	UnitID                      int    `xorm:"unit_id"`
	Exp                         int    `xorm:"exp"`
	NextExp                     int    `xorm:"next_exp"`
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
	FavoriteFlag                int    `xorm:"favorite_flag"`
	DisplayRank                 int    `xorm:"display_rank"`
	IsRankMax                   int    `xorm:"is_rank_max"`
	IsLoveMax                   int    `xorm:"is_love_max"`
	IsLevelMax                  int    `xorm:"is_level_max"`
	IsSigned                    int    `xorm:"is_signed"`
	IsSkillLevelMax             int    `xorm:"is_skill_level_max"`
	IsRemovableSkillCapacityMax int    `xorm:"is_removable_skill_capacity_max"`
	InsertDate                  string `xorm:"insert_date"`
}

func (UserUnit) TableName() string {
	return "user_unit"
}
