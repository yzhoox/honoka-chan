package user

type Users struct {
	ID            int    `xorm:"id pk autoincr"`
	Phone         string `xorm:"phone"`
	Password      string `xorm:"password"`
	Autokey       string `xorm:"autokey"`
	Ticket        string `xorm:"ticket"`
	UserID        int    `xorm:"user_id"`
	LastLoginTime int64  `xorm:"last_login_time"`
}

func (Users) TableName() string {
	return "users"
}
