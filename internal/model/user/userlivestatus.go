package usermodel

type UserLiveStatus struct {
	ID               int   `xorm:"id pk autoincr"`
	UserID           int   `xorm:"user_id unique(user_live_status_pair) index"`
	LiveDifficultyID int   `xorm:"live_difficulty_id unique(user_live_status_pair) index"`
	HiScore          int   `xorm:"hi_score"`
	HiComboCount     int   `xorm:"hi_combo_count"`
	ClearCnt         int   `xorm:"clear_cnt"`
	InsertDate       int64 `xorm:"insert_date"`
	UpdateDate       int64 `xorm:"update_date"`
}

func (UserLiveStatus) TableName() string {
	return "user_live_status"
}
