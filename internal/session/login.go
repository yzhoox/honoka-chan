package session

import loginmodel "honoka-chan/internal/model/login"

func (ss *Session) GetAuthKey(token string) (bool, *loginmodel.AuthKey) {
	authKeyData := loginmodel.AuthKey{}
	has, err := ss.UserEng.Table(new(loginmodel.AuthKey)).
		Where("authorize_token = ?", token).Get(&authKeyData)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &authKeyData
}

func (ss *Session) SetAuthKey(authKey *loginmodel.AuthKey) {
	_, err := ss.UserEng.Insert(authKey)
	if ss.CheckErr(err) {
		return
	}
}
