package usermodel

type UserLiveInProgress struct {
	ID     int `xorm:"id pk autoincr"`
	DeckID int `xorm:"deck_id"`
	UserID int `xorm:"user_id"`
}

func (UserLiveInProgress) TableName() string {
	return "user_live_in_progress"
}
