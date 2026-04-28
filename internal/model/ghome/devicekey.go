package ghomemodel

type DeviceKey struct {
	ID       int    `xorm:"id pk autoincr"`
	DeviceID string `xorm:"device_id"`
	RandKey  string `xorm:"rand_key"`
}

func (DeviceKey) TableName() string {
	return "device_key"
}
