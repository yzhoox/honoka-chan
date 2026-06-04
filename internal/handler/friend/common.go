package friend

import (
	"errors"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/session"
	"strconv"
	"strings"
)

func resolveActualFriendUserID(ss *session.Session, requestedUserID int) (int, error) {
	if requestedUserID <= 0 {
		return 0, errors.New("invalid user_id")
	}

	pref := usermodel.UserPref{}
	has, err := ss.UserEng.Table(new(usermodel.UserPref)).
		Where("invite_code = ?", strconv.Itoa(requestedUserID)).
		Cols("user_id").
		Get(&pref)
	if err != nil {
		return 0, err
	}
	if has && pref.UserID > 0 {
		return pref.UserID, nil
	}

	user := usermodel.Users{}
	has, err = ss.UserEng.Table(new(usermodel.Users)).
		Where("user_id = ?", requestedUserID).
		Cols("user_id").
		Get(&user)
	if err != nil {
		return 0, err
	}
	if !has || user.UserID <= 0 {
		return 0, errors.New("user not found")
	}

	return user.UserID, nil
}

func areUsersFriends(ss *session.Session, userID, friendUserID int) (bool, error) {
	forward, err := ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", userID).
		Where("friend_user_id = ?", friendUserID).
		Where("status = ?", usermodel.FriendStatusApproved).
		Exist()
	if err != nil {
		return false, err
	}
	if !forward {
		return false, nil
	}

	return ss.UserEng.Table(new(usermodel.UserFriend)).
		Where("user_id = ?", friendUserID).
		Where("friend_user_id = ?", userID).
		Where("status = ?", usermodel.FriendStatusApproved).
		Exist()
}

func resolveUserIDByInviteCode(ss *session.Session, inviteCode string) (int, error) {
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, inviteCode)
	if digits == "" {
		return 0, errors.New("invalid invite_code")
	}

	pref := usermodel.UserPref{}
	has, err := ss.UserEng.Table(new(usermodel.UserPref)).
		Where("invite_code = ?", digits).
		Cols("user_id").
		Get(&pref)
	if err != nil {
		return 0, err
	}
	if !has || pref.UserID <= 0 {
		return 0, errors.New("user not found")
	}

	return pref.UserID, nil
}
