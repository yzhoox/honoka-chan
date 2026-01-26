package usermodel

type UserKey struct {
	ID     int `xorm:"id pk autoincr"`
	UserID int `xorm:"user_id"`
	Key    int `xorm:"key"`
}

func (UserKey) TableName() string {
	return "user_key"
}
