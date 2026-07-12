package loginmodel

type AuthKey struct {
	ID             int    `xorm:"id pk autoincr"`
	AuthorizeToken string `xorm:"authorize_token"`
	UserID         int    `xorm:"user_id index"`
	ClientToken    string `xorm:"client_token"`
	ServerToken    string `xorm:"server_token"`
	InsertDate     string `xorm:"insert_date"`
}

func (AuthKey) TableName() string {
	return "auth_key"
}
