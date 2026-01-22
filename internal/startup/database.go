package startup

import (
	"honoka-chan/internal/model/user"
	"honoka-chan/pkg/db"
)

func CreateTable() {
	// db.UserEng.ShowSQL(true)
	db.UserEng.Sync2(new(user.Users))
	db.UserEng.Sync2(new(user.UserKey))
	db.UserEng.Sync2(new(user.UserPref))
	db.UserEng.Sync2(new(user.UserAccessoryWear))
	db.UserEng.Sync2(new(user.UserDeck))
	db.UserEng.Sync2(new(user.UserUnit))
	db.UserEng.Sync2(new(user.UserDeckUnit))
	db.UserEng.Sync2(new(user.UserUnitSkillEquip))
}
