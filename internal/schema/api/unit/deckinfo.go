package unitapischema

type UserDeckData struct {
	ID         int    `xorm:"id pk autoincr"`
	DeckID     int    `xorm:"deck_id"`
	MainFlag   int    `xorm:"main_flag"`
	DeckName   string `xorm:"deck_name"`
	UserID     int    `xorm:"user_id"`
	InsertDate int64  `xorm:"insert_date"`
}

type UnitDeckData struct {
	ID               int   `xorm:"id pk autoincr" json:"-"`
	UserDeckID       int   `xorm:"user_deck_id" json:"-"`
	UnitOwningUserID int   `xorm:"unit_owning_user_id" json:"unit_owning_user_id"`
	UnitID           int   `xorm:"unit_id" json:"unit_id"`
	Position         int   `xorm:"position" json:"position"`
	Level            int   `xorm:"level" json:"level"`
	LevelLimitID     int   `xorm:"level_limit_id" json:"level_limit_id"`
	DisplayRank      int   `xorm:"display_rank" json:"display_rank"`
	Love             int   `xorm:"love" json:"love"`
	UnitSkillLevel   int   `xorm:"unit_skill_level" json:"unit_skill_level"`
	IsRankMax        bool  `xorm:"is_rank_max" json:"is_rank_max"`
	IsLoveMax        bool  `xorm:"is_love_max" json:"is_love_max"`
	IsLevelMax       bool  `xorm:"is_level_max" json:"is_level_max"`
	IsSigned         bool  `xorm:"is_signed" json:"is_signed"`
	BeforeLove       int   `xorm:"before_love" json:"before_love"`
	MaxLove          int   `xorm:"max_love" json:"max_love"`
	InsertDate       int64 `xorm:"insert_date" json:"-"`
}

type UnitOwningUserIds struct {
	Position         int `json:"position"`
	UnitOwningUserID int `json:"unit_owning_user_id"`
}

type DeckInfoData struct {
	UnitDeckID        int                 `json:"unit_deck_id"`
	MainFlag          bool                `json:"main_flag"`
	DeckName          string              `json:"deck_name"`
	UnitOwningUserIds []UnitOwningUserIds `json:"unit_owning_user_ids"`
}

type DeckInfoResp struct {
	Result     []DeckInfoData `json:"result"`
	Status     int            `json:"status"`
	CommandNum bool           `json:"commandNum"`
	TimeStamp  int64          `json:"timeStamp"`
}
