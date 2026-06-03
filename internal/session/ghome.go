package session

import (
	"errors"
	ghomemodel "honoka-chan/internal/model/ghome"
)

func (ss *Session) GetDeviceID() string {
	return ss.deviceID
}

func (ss *Session) GetRandKey() ([]byte, error) {
	deviceKey := ghomemodel.DeviceKey{}
	has, err := ss.UserEng.Table(new(ghomemodel.DeviceKey)).
		Where("device_id = ?", ss.GetDeviceID()).Get(&deviceKey)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return []byte(deviceKey.RandKey), nil
}
func (ss *Session) SetRandKey(key string) {
	var err error
	randKey, err := ss.GetRandKey()
	if ss.CheckErr(err) {
		return
	}
	if randKey == nil {
		_, err = ss.UserEng.Insert(&ghomemodel.DeviceKey{
			DeviceID: ss.GetDeviceID(),
			RandKey:  key,
		})
	} else {
		_, err = ss.UserEng.Table(new(ghomemodel.DeviceKey)).
			Where("device_id = ?", ss.GetDeviceID()).Update(&ghomemodel.DeviceKey{RandKey: key})
	}
	if ss.CheckErr(err) {
		return
	}
}

func (ss *Session) Get3DESRandKey() ([]byte, error) {
	randKey, err := ss.GetRandKey()
	if err != nil {
		return nil, err
	}
	if len(randKey) < 24 {
		return nil, errors.New("invalid rand key")
	}
	return randKey[:24], nil
}
