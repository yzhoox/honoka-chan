package usermodel

import (
	"strconv"
	"strings"
	"time"

	"xorm.io/xorm"
)

const CurrentUserPrefProfileVersion = 2

const (
	DefaultAutoUserName      = "音乃木坂学生"
	DefaultAutoUserDesc      = "你好。"
	DefaultUserLevel         = 1028
	DefaultUserExp           = 1089696
	DefaultUserNextExp       = 1207185
	DefaultUserGameCoin      = 112124104
	DefaultUserSnsCoin       = 10000
	DefaultUserEnergyMax     = 417
	DefaultUserOverMaxEnergy = 0
	DefaultBirthMonth        = 10
	DefaultBirthDay          = 18
)

type UserPref struct {
	ID               int    `xorm:"id pk autoincr"`
	UserID           int    `xorm:"user_id"`
	AwardID          int    `xorm:"award_id"`
	BackgroundID     int    `xorm:"background_id"`
	UnitOwningUserID int    `xorm:"unit_owning_user_id"`
	UserName         string `xorm:"user_name"`
	UserLevel        int    `xorm:"user_level"`
	UserDesc         string `xorm:"user_desc"`
	InviteCode       string `xorm:"invite_code"`
	UserExp          int    `xorm:"user_exp"`
	NextExp          int    `xorm:"next_exp"`
	GameCoin         int    `xorm:"game_coin"`
	SnsCoin          int    `xorm:"sns_coin"`
	EnergyMax        int    `xorm:"energy_max"`
	OverMaxEnergy    int    `xorm:"over_max_energy"`
	BirthMonth       int    `xorm:"birth_month"`
	BirthDay         int    `xorm:"birth_day"`
	ForceRelogin     bool   `xorm:"force_relogin"`
	ProfileVersion   int    `xorm:"profile_version"`
	UpdateTime       int64  `xorm:"update_time"`
}

func (u *UserPref) ApplyProfileDefaults() {
	if u.UserLevel <= 0 {
		u.UserLevel = DefaultUserLevel
	}
	if strings.TrimSpace(u.InviteCode) == "" && u.UserID > 0 {
		u.InviteCode = strconv.Itoa(u.UserID)
	}
	if u.UserExp < 0 || u.UserExp == 0 {
		u.UserExp = DefaultUserExp
	}
	if u.NextExp < 0 || u.NextExp == 0 {
		u.NextExp = DefaultUserNextExp
	}
	if u.GameCoin < 0 || u.GameCoin == 0 {
		u.GameCoin = DefaultUserGameCoin
	}
	if u.SnsCoin < 0 || u.SnsCoin == 0 {
		u.SnsCoin = DefaultUserSnsCoin
	}
	if u.EnergyMax <= 0 {
		u.EnergyMax = DefaultUserEnergyMax
	}
	if u.OverMaxEnergy < 0 {
		u.OverMaxEnergy = DefaultUserOverMaxEnergy
	}
	if !IsValidBirthDate(u.BirthMonth, u.BirthDay) {
		u.BirthMonth = DefaultBirthMonth
		u.BirthDay = DefaultBirthDay
	}
	u.ProfileVersion = CurrentUserPrefProfileVersion
}

func (u *UserPref) ResetProfileDefaults() {
	u.UserName = DefaultAutoUserName
	u.UserDesc = DefaultAutoUserDesc
	u.UserLevel = DefaultUserLevel
	u.UserExp = DefaultUserExp
	u.NextExp = DefaultUserNextExp
	u.GameCoin = DefaultUserGameCoin
	u.SnsCoin = DefaultUserSnsCoin
	u.EnergyMax = DefaultUserEnergyMax
	u.OverMaxEnergy = DefaultUserOverMaxEnergy
	u.BirthMonth = DefaultBirthMonth
	u.BirthDay = DefaultBirthDay
	if u.UserID > 0 {
		u.InviteCode = strconv.Itoa(u.UserID)
	}
	u.ProfileVersion = CurrentUserPrefProfileVersion
}

func (u UserPref) NeedsProfileMigration() bool {
	return u.ProfileVersion < CurrentUserPrefProfileVersion
}

func (u UserPref) EffectiveEnergyMax() int {
	if u.EnergyMax > 0 {
		return u.EnergyMax
	}
	return DefaultUserEnergyMax
}

func (u UserPref) EffectiveCurrentEnergy() int {
	if u.OverMaxEnergy > 0 {
		return u.OverMaxEnergy
	}
	return u.EffectiveEnergyMax()
}

func (u UserPref) EffectiveBirthMonth() int {
	if IsValidBirthDate(u.BirthMonth, u.BirthDay) {
		return u.BirthMonth
	}
	return DefaultBirthMonth
}

func (u UserPref) EffectiveBirthDay() int {
	if IsValidBirthDate(u.BirthMonth, u.BirthDay) {
		return u.BirthDay
	}
	return DefaultBirthDay
}

func (u UserPref) HasBirthDate() bool {
	return IsValidBirthDate(u.BirthMonth, u.BirthDay)
}

func (u UserPref) IsBirthdayToday(now time.Time) bool {
	return int(now.Month()) == u.EffectiveBirthMonth() && now.Day() == u.EffectiveBirthDay()
}

func IsValidBirthDate(month, day int) bool {
	if month < 1 || month > 12 || day < 1 {
		return false
	}
	maxDay := time.Date(2000, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	return day <= maxDay
}

func UserPrefProfileColumns() []string {
	return []string{
		"user_level",
		"invite_code",
		"user_exp",
		"next_exp",
		"game_coin",
		"sns_coin",
		"energy_max",
		"over_max_energy",
		"birth_month",
		"birth_day",
		"force_relogin",
		"profile_version",
	}
}

func MarkUserForceRelogin(exec xorm.Interface, userID int) error {
	_, err := exec.Table(new(UserPref)).
		Where("user_id = ?", userID).
		Cols("force_relogin").
		Update(&UserPref{ForceRelogin: true})
	return err
}

func ClearUserForceRelogin(exec xorm.Interface, userID int) error {
	_, err := exec.Table(new(UserPref)).
		Where("user_id = ?", userID).
		Cols("force_relogin").
		Update(&UserPref{ForceRelogin: false})
	return err
}

func (UserPref) TableName() string {
	return "user_pref"
}
