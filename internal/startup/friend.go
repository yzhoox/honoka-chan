package startup

import (
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/pkg/db"
	"log"
)

func EnsureDefaultFriends() {
	session := db.UserEng.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		log.Fatalln("初始化默认好友失败:", err.Error())
	}

	userIDs := []int{}
	err := session.Table(new(usermodel.Users)).
		Cols("user_id").
		Find(&userIDs)
	if err != nil {
		session.Rollback()
		log.Fatalln("初始化默认好友失败:", err.Error())
	}

	for _, userID := range userIDs {
		err = usermodel.EnsureDefaultFriendship(session, userID)
		if err != nil {
			session.Rollback()
			log.Fatalln("初始化默认好友失败:", err.Error())
		}
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		log.Fatalln("初始化默认好友失败:", err.Error())
	}
}
