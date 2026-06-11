package startup

import (
	"fmt"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/pkg/db"
)

func EnsureDefaultFriends() error {
	session := db.UserEng.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return fmt.Errorf("初始化默认好友失败: %w", err)
	}

	userIDs := []int{}
	err := session.Table(new(usermodel.Users)).
		Cols("user_id").
		Find(&userIDs)
	if err != nil {
		session.Rollback()
		return fmt.Errorf("初始化默认好友失败: %w", err)
	}

	for _, userID := range userIDs {
		err = usermodel.EnsureDefaultFriendship(session, userID)
		if err != nil {
			session.Rollback()
			return fmt.Errorf("初始化默认好友失败: %w", err)
		}
	}

	if err := session.Commit(); err != nil {
		session.Rollback()
		return fmt.Errorf("初始化默认好友失败: %w", err)
	}
	return nil
}
