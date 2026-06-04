package usermodel

const GreetingPageSize = 40

type UserGreet struct {
	NoticeID            int    `xorm:"notice_id pk autoincr"`
	AffectorID          int    `xorm:"affector_id index"`
	ReceiverID          int    `xorm:"receiver_id index"`
	Message             string `xorm:"message text"`
	Reply               bool   `xorm:"reply"`
	Readed              bool   `xorm:"readed"`
	DeletedFromAffector bool   `xorm:"deleted_from_affector"`
	DeletedFromReceiver bool   `xorm:"deleted_from_receiver"`
	InsertDate          int64  `xorm:"insert_date index"`
}

func (UserGreet) TableName() string {
	return "user_greet"
}
