package session

import (
	accessorymodel "honoka-chan/internal/model/accessory"
	usermodel "honoka-chan/internal/model/user"
	liverecordschema "honoka-chan/internal/schema/liverecord"
)

func (ss *Session) GetUserAccessoryWearByUnitOwningUserID(unitOwningUserID int) (bool, *usermodel.UserAccessoryWear) {
	wearData := usermodel.UserAccessoryWear{}
	has, err := ss.UserEng.Table(new(usermodel.UserAccessoryWear)).
		Where("user_id = ? AND unit_owning_user_id = ?", ss.UserID, unitOwningUserID).Get(&wearData)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &wearData
}

func (ss *Session) GetUserAccessoryByAccessoryOwningUserID(accessoryOwningUserID int) (bool, *usermodel.UserAccessory) {
	return ss.GetUserAccessoryByAccessoryOwningUserIDForUser(ss.UserID, accessoryOwningUserID)
}

func (ss *Session) GetUserAccessoryByAccessoryOwningUserIDForUser(userID, accessoryOwningUserID int) (bool, *usermodel.UserAccessory) {
	accessoryData := usermodel.UserAccessory{}
	has, err := ss.UserEng.Table(new(usermodel.UserAccessory)).
		Where("user_id = ? AND accessory_owning_user_id = ?", userID, accessoryOwningUserID).
		Get(&accessoryData)
	if ss.CheckErr(err) {
		return false, nil
	}

	return has, &accessoryData
}

func (ss *Session) GetAccessoryByAccessoryOwningUserID(accessoryOwningUserID int) (bool, *accessorymodel.Accessory) {
	return ss.GetAccessoryByAccessoryOwningUserIDForUser(ss.UserID, accessoryOwningUserID)
}

func (ss *Session) GetAccessoryByAccessoryOwningUserIDForUser(userID, accessoryOwningUserID int) (bool, *accessorymodel.Accessory) {
	has, userAccessoryData := ss.GetUserAccessoryByAccessoryOwningUserIDForUser(userID, accessoryOwningUserID)
	if !has {
		return false, nil
	}

	type accessoryWithUserExp struct {
		accessorymodel.Accessory `xorm:"extends"`
		Exp                      int `xorm:"-"`
	}

	accessoryData := accessoryWithUserExp{}
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

	return has, &accessoryData.Accessory
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
