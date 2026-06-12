package session

import (
	accessorymodel "honoka-chan/internal/model/accessory"
	usermodel "honoka-chan/internal/model/user"
	liverecordschema "honoka-chan/internal/schema/liverecord"
)

func (ss *Session) GetUserAccessoryWearByUnitOwningUserID(unitOwningUserID int) (bool, *usermodel.UserAccessoryWear) {
	wearData := usermodel.UserAccessoryWear{}
	has, err := ss.UserEng.Table(new(usermodel.UserAccessoryWear)).
		Where("unit_owning_user_id = ?", unitOwningUserID).Get(&wearData)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &wearData
}

func (ss *Session) GetUserAccessoryByAccessoryOwningUserID(accessoryOwningUserID int) (bool, *usermodel.UserAccessory) {
	accessoryData := usermodel.UserAccessory{}
	has, err := ss.UserEng.Table(new(usermodel.UserAccessory)).
		Where("accessory_owning_user_id = ?", accessoryOwningUserID).
		Get(&accessoryData)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &accessoryData
}

func (ss *Session) GetAccessoryByAccessoryOwningUserID(accessoryOwningUserID int) (bool, *accessorymodel.Accessory) {
	has, userAccessoryData := ss.GetUserAccessoryByAccessoryOwningUserID(accessoryOwningUserID)
	if !has {
		return false, nil
	}

	accessoryData := accessorymodel.Accessory{}
	has, err := ss.MainEng.Table("accessory_m").
		Where("accessory_id = ?", userAccessoryData.AccessoryID).
		Get(&accessoryData)
	if ss.CheckErr(err) {
		return false, nil
	}
	if !has {
		return false, nil
	}

	accessoryData.Exp = userAccessoryData.Exp

	return has, &accessoryData
}

func (ss *Session) GetUserAccessoryInfoByUnitOwningUserID(unitOwningUserID int) (bool, *liverecordschema.AccessoryInfo) {
	has, wearData := ss.GetUserAccessoryWearByUnitOwningUserID(unitOwningUserID)
	if !has {
		return false, nil
	}

	has, accessoryData := ss.GetAccessoryByAccessoryOwningUserID(wearData.AccessoryOwningUserID)
	if !has {
		return false, nil
	}

	return has, &liverecordschema.AccessoryInfo{
		AccessoryOwningUserID: wearData.AccessoryOwningUserID,
		AccessoryID:           accessoryData.AccessoryID,
		Exp:                   accessoryData.Exp,
		NextExp:               0,
		Level:                 accessoryData.MaxLevel,
		MaxLevel:              accessoryData.MaxLevel,
		RankUpCount:           accessoryData.MaxLevel - accessoryData.DefaultMaxLevel,
		FavoriteFlag:          true,
	}
}
