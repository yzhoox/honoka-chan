package db

import (
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var (
	MainDb  = "assets/main.db"
	UserDb  = "assets/data.db"
	MainEng *xorm.Engine
	UserEng *xorm.Engine
)

func init() {
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
