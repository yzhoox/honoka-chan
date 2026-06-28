package accessorymodel

type Accessory struct {
	AccessoryID     int `xorm:"accessory_id pk"`
	Exp             int `xorm:"-"`
	SmileMax        int `xorm:"smile_max"`
	PureMax         int `xorm:"pure_max"`
	CoolMax         int `xorm:"cool_max"`
	DefaultMaxLevel int `xorm:"default_max_level"`
	MaxLevel        int `xorm:"max_level"`
}
