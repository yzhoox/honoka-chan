package usermodel

type UserPref struct {
	ID               int    `xorm:"id pk autoincr"`
	UserID           int    `xorm:"user_id"`
	AwardID          int    `xorm:"award_id"`
	BackgroundID     int    `xorm:"background_id"`
	UnitOwningUserID int    `xorm:"unit_owning_user_id"`
	UserName         string `xorm:"user_name"`
	UserLevel        int    `xorm:"user_level"`
	UserDesc         string `xorm:"user_desc"`
	InviteCode       string `xorm:"invite_code"`
	UpdateTime       int64  `xorm:"update_time"`
}

func (UserPref) TableName() string {
	return "user_pref"
}
