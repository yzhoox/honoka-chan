package usermodel

type UserAccessoryWear struct {
	ID                    int `xorm:"id pk autoincr"`
	AccessoryOwningUserID int `xorm:"accessory_owning_user_id"`
	UnitOwningUserID      int `xorm:"unit_owning_user_id"`
	UserID                int `xorm:"user_id"`
}

func (u *UserAccessoryWear) TableName() string {
	return "user_accessory_wear"
}
