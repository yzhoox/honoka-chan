package usermodel

type UserUnitData struct {
	UnitOwningUserID int   `xorm:"unit_owning_user_id pk autoincr"`
	UnitID           int   `xorm:"unit_id"`
	FavoriteFlag     bool  `xorm:"favorite_flag"`
	DisplayRank      int   `xorm:"display_rank"`
	UserID           int   `xorm:"user_id"`
	InsertDate       int64 `xorm:"insert_date"`
}

func (UserUnitData) TableName() string {
	return "user_unit_data"
}
