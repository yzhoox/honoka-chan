package usermodel

type UserLiveRandom struct {
	UserID           int    `xorm:"user_id pk"`
	Attribute        int    `xorm:"attribute pk"`
	Difficulty       int    `xorm:"difficulty pk"`
	MemberCategory   int    `xorm:"member_category pk"`
	Token            string `xorm:"token"`
	LiveDifficultyID int    `xorm:"live_difficulty_id"`
	InProgress       bool   `xorm:"in_progress"`
}

func (UserLiveRandom) TableName() string {
	return "user_live_random"
}
