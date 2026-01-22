package db

import (
	"fmt"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var (
	Ldb *LdbInstance

	MainDb  = "assets/main.db"
	UserDb  = "assets/data.db"
	MainEng *xorm.Engine
	UserEng *xorm.Engine
)

func init() {
	Ldb = GetLdbInstance()

	eng, err := xorm.NewEngine("sqlite", MainDb)
	if err != nil {
		panic(err)
	}
	err = eng.Ping()
	if err != nil {
		panic(err)
	}
	eng.SetMaxOpenConns(10)
	eng.SetMaxIdleConns(5)
	MainEng = eng

	eng, err = xorm.NewEngine("sqlite", UserDb)
	if err != nil {
		panic(err)
	}
	err = eng.Ping()
	if err != nil {
		panic(err)
	}
	eng.SetMaxOpenConns(10)
	eng.SetMaxIdleConns(5)
	UserEng = eng
}

func MatchTokenUid(token, uid string) bool {
	res, err := Ldb.Get([]byte(uid))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return string(res) == token
}
