package session

import loginmodel "honoka-chan/internal/model/login"

func (ss *Session) GetAuthKey(token string) (bool, *loginmodel.AuthKey, error) {
	authKeyData := loginmodel.AuthKey{}
	has, err := ss.UserEng.Table(new(loginmodel.AuthKey)).
		Where("authorize_token = ?", token).Get(&authKeyData)
	if err != nil {
		return false, nil, err
	}

	return has, &authKeyData, nil
}

func (ss *Session) SetAuthKey(authKey *loginmodel.AuthKey) {
	_, err := ss.UserEng.Insert(authKey)
	if ss.CheckErr(err) {
		return
	}
}

func (ss *Session) IsAuthorizeTokenForUser(token string, userID int) (bool, error) {
	if token == "" || userID <= 0 {
		return false, nil
	}

	return ss.UserEng.Table(new(loginmodel.AuthKey)).
		Where("authorize_token = ? AND user_id = ?", token, userID).
		Exist()
}
