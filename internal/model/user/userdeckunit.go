package usermodel

type UserDeckUnit struct {
	ID               int   `xorm:"id pk autoincr"`
	UserDeckID       int   `xorm:"user_deck_id"`
	UnitOwningUserID int   `xorm:"unit_owning_user_id"`
	UnitID           int   `xorm:"unit_id"`
	Position         int   `xorm:"position"`
	Level            int   `xorm:"level"`
	LevelLimitID     int   `xorm:"level_limit_id"`
	DisplayRank      int   `xorm:"display_rank"`
	Love             int   `xorm:"love"`
	UnitSkillLevel   int   `xorm:"unit_skill_level"`
	IsRankMax        bool  `xorm:"is_rank_max"`
	IsLoveMax        bool  `xorm:"is_love_max"`
	IsLevelMax       bool  `xorm:"is_level_max"`
	IsSigned         bool  `xorm:"is_signed"`
	BeforeLove       int   `xorm:"before_love"`
	MaxLove          int   `xorm:"max_love"`
	UserID           int   `xorm:"user_id"`
	InsertDate       int64 `xorm:"insert_date"`
}

func (UserDeckUnit) TableName() string {
	return "user_deck_unit"
}
