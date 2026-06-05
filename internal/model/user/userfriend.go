package usermodel

import (
	"time"

	"xorm.io/xorm"
)

const (
	DefaultSystemPhone = "1"

	FriendStatusApproved = iota
	FriendStatusPending
	FriendStatusAwaitingApproval

	FriendListPageSize = 40
)

const (
	ClientFriendStatusNone = iota
	ClientFriendStatusApproved
	ClientFriendStatusPending
	ClientFriendStatusAwaitingApproval
)

type UserFriend struct {
	ID           int   `xorm:"id pk autoincr"`
	UserID       int   `xorm:"user_id unique(friend_pair) index"`
	FriendUserID int   `xorm:"friend_user_id unique(friend_pair) index"`
	Status       int   `xorm:"status index"`
	IsNew        bool  `xorm:"is_new"`
	InsertDate   int64 `xorm:"insert_date index"`
	UpdateDate   int64 `xorm:"update_date"`
}

func (UserFriend) TableName() string {
	return "user_friend"
}

func EnsureFriendLink(dbSession *xorm.Session, userID, friendUserID, status int, isNew bool) error {
	if dbSession == nil || userID <= 0 || friendUserID <= 0 || userID == friendUserID {
		return nil
	}

	now := time.Now().Unix()
	link := UserFriend{}
	has, err := dbSession.Table(new(UserFriend)).
		Where("user_id = ? AND friend_user_id = ?", userID, friendUserID).
		Get(&link)
	if err != nil {
		return err
	}

	if has {
		update := UserFriend{
			Status:     status,
			IsNew:      isNew,
			UpdateDate: now,
		}
		if link.InsertDate <= 0 {
			update.InsertDate = now
		}

		_, err = dbSession.Table(new(UserFriend)).
			Where("user_id = ? AND friend_user_id = ?", userID, friendUserID).
			Cols("status", "is_new", "update_date", "insert_date").
			Update(&update)
		return err
	}

	_, err = dbSession.Table(new(UserFriend)).Insert(&UserFriend{
		UserID:       userID,
		FriendUserID: friendUserID,
		Status:       status,
		IsNew:        isNew,
		InsertDate:   now,
		UpdateDate:   now,
	})
	return err
}

func EnsureMutualFriend(dbSession *xorm.Session, userID, friendUserID, status int) error {
	if err := EnsureFriendLink(dbSession, userID, friendUserID, status, false); err != nil {
		return err
	}
	return EnsureFriendLink(dbSession, friendUserID, userID, status, false)
}

func EnsureMutualFriendWithIsNew(dbSession *xorm.Session, userID, friendUserID, status int, userIsNew, friendIsNew bool) error {
	if err := EnsureFriendLink(dbSession, userID, friendUserID, status, userIsNew); err != nil {
		return err
	}
	return EnsureFriendLink(dbSession, friendUserID, userID, status, friendIsNew)
}

func ResolveClientFriendStatus(dbSession *xorm.Session, userID, targetUserID int) (int, error) {
	if dbSession == nil || userID <= 0 || targetUserID <= 0 || userID == targetUserID {
		return ClientFriendStatusNone, nil
	}

	link := UserFriend{}
	has, err := dbSession.Table(new(UserFriend)).
		Where("user_id = ? AND friend_user_id = ?", userID, targetUserID).
		Cols("status").
		Get(&link)
	if err != nil {
		return ClientFriendStatusNone, err
	}
	if !has {
		return ClientFriendStatusNone, nil
	}

	switch link.Status {
	case FriendStatusApproved:
		isMutual, err := dbSession.Table(new(UserFriend)).
			Where("user_id = ? AND friend_user_id = ?", targetUserID, userID).
			Where("status = ?", FriendStatusApproved).
			Exist()
		if err != nil {
			return ClientFriendStatusNone, err
		}
		if isMutual {
			return ClientFriendStatusApproved, nil
		}
	case FriendStatusPending:
		return ClientFriendStatusPending, nil
	case FriendStatusAwaitingApproval:
		return ClientFriendStatusAwaitingApproval, nil
	}

	return ClientFriendStatusNone, nil
}

func EnsureDefaultFriendship(dbSession *xorm.Session, userID int) error {
	if dbSession == nil || userID <= 0 {
		return nil
	}

	defaultUser := Users{}
	has, err := dbSession.Table(new(Users)).
		Where("phone = ?", DefaultSystemPhone).
		Cols("user_id").
		Get(&defaultUser)
	if err != nil {
		return err
	}
	if !has || defaultUser.UserID <= 0 || defaultUser.UserID == userID {
		return nil
	}

	return EnsureMutualFriend(dbSession, userID, defaultUser.UserID, FriendStatusApproved)
}
