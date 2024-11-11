package db

import "fmt"

var (
	DB *Instance
)

func init() {
	DB = GetInstance()
}

func MatchTokenUid(token, uid string) bool {
	res, err := DB.Get([]byte(uid))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return string(res) == token
}
