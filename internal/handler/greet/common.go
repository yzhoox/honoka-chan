package greet

import (
	"errors"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/session"
	"strconv"
)

func resolveGreetingUserID(ss *session.Session, requestedUserID int) (int, error) {
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
