package db

import (
	"fmt"
	"honoka-chan/pkg/utils"
	"os"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var (
	DB *Instance

	ExampleDb = "assets/data.example.db"
	MainDb    = "assets/main.db"
	UserDb    = "assets/data.db"
	MainEng   *xorm.Engine
	UserEng   *xorm.Engine
)

func init() {
	DB = GetInstance()

	_, err := os.Stat(UserDb)
	if err != nil {
		utils.WriteAllText(UserDb, utils.ReadAllText(ExampleDb))
	}

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
	res, err := DB.Get([]byte(uid))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return string(res) == token
}
