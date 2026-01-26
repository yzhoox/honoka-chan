package loginmodel

type AuthKey struct {
	ID             int    `xorm:"id pk autoincr"`
	AuthorizeToken string `xorm:"authorize_token"`
	ClientToken    string `xorm:"client_token"`
	ServerToken    string `xorm:"server_token"`
	InsertDate     string `xorm:"insert_date"`
}
