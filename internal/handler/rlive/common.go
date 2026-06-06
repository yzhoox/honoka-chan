package rlive

import (
	"errors"
	usermodel "honoka-chan/internal/model/user"
	"honoka-chan/internal/session"
)

func getRandomLiveByToken(ss *session.Session, token string) (usermodel.UserLiveRandom, error) {
	if token == "" {
		return usermodel.UserLiveRandom{}, errors.New("random live token is required")
	}

	row := usermodel.UserLiveRandom{}
	has, err := ss.UserEng.Table(new(usermodel.UserLiveRandom)).
		Where("user_id = ? AND token = ?", ss.UserID, token).
		Get(&row)
	if err != nil {
		return usermodel.UserLiveRandom{}, err
	}
	if !has {
		return usermodel.UserLiveRandom{}, errors.New("random live session not found")
	}

	return row, nil
}

func deleteRandomLiveByToken(ss *session.Session, token string) {
	if token == "" || ss.UserEng == nil {
		return
	}

	_, err := ss.UserEng.Table(new(usermodel.UserLiveRandom)).
		Where("user_id = ? AND token = ?", ss.UserID, token).
		Delete()
	if ss.CheckErr(err) {
		return
	}
}
