package usermodel

type UserDeck struct {
	ID         int    `xorm:"id pk autoincr"`
	DeckID     int    `xorm:"deck_id"`
	DeckName   string `xorm:"deck_name"`
	MainFlag   int    `xorm:"main_flag"`
	UserID     int    `xorm:"user_id"`
	InsertDate int64  `xorm:"insert_date"`
}

func (UserDeck) TableName() string {
	return "user_deck"
}
