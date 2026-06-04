package usermodel

import (
	"strconv"
	"strings"
)

const CurrentUserPrefProfileVersion = 1

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
		"profile_version",
	}
}

func (UserPref) TableName() string {
	return "user_pref"
}
