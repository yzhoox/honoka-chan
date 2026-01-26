package session

import (
	ghomemodel "honoka-chan/internal/model/ghome"
)

func (ss *Session) GetDeviceID() string {
	return ss.Ctx.Request.Header.Get("X-DEVICEID")
}

func (ss *Session) GetRandKey() []byte {
	deviceKey := ghomemodel.DeviceKey{}
	has, err := ss.UserEng.Table(new(ghomemodel.DeviceKey)).
		Where("device_id = ?", ss.GetDeviceID()).Get(&deviceKey)
	if ss.CheckErr(err) {
		return nil
	}

	if !has {
		return nil
	}

	return []byte(deviceKey.RandKey)
}
func (ss *Session) SetRandKey(key string) {
	var err error
	if ss.GetRandKey() == nil {
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
