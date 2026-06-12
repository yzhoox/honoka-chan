package usermodel

import (
	"fmt"
	"sort"

	"xorm.io/xorm"
)

type UserAccessory struct {
	AccessoryOwningUserID int `xorm:"accessory_owning_user_id pk autoincr"`
	UserID                int `xorm:"user_id index"`
	AccessoryID           int `xorm:"accessory_id index"`
	Exp                   int `xorm:"exp"`
}

func (UserAccessory) TableName() string {
	return "user_accessory"
}

type accessoryTemplateRow struct {
	AccessoryID int `xorm:"accessory_id"`
}

type userAccessoryCountRow struct {
	AccessoryID int `xorm:"accessory_id"`
	Count       int `xorm:"cnt"`
}

func LoadAccessoryTemplateIDs(mainSession xorm.Interface) ([]int, error) {
	rows := []accessoryTemplateRow{}
	err := mainSession.Table("accessory_m").
		Where("is_material = ? AND rarity = ?", 0, 4).
		Cols("accessory_id").
		OrderBy("accessory_id ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.AccessoryID)
	}
	return ids, nil
}

func EnsureUserAccessories(userSession *xorm.Session, mainSession xorm.Interface, userID int, minCount map[int]int) error {
	accessoryIDs, err := LoadAccessoryTemplateIDs(mainSession)
	if err != nil {
		return err
	}

	existingRows := []userAccessoryCountRow{}
	err = userSession.Table(new(UserAccessory)).
		Select("accessory_id, COUNT(*) AS cnt").
		Where("user_id = ?", userID).
		GroupBy("accessory_id").
		Find(&existingRows)
	if err != nil {
		return err
	}

	existingCount := make(map[int]int, len(existingRows))
	for _, row := range existingRows {
		existingCount[row.AccessoryID] = row.Count
	}

	for _, accessoryID := range accessoryIDs {
		targetCount := 1
		if count, ok := minCount[accessoryID]; ok && count > targetCount {
			targetCount = count
		}

		for existingCount[accessoryID] < targetCount {
			_, err = userSession.Insert(&UserAccessory{
				UserID:      userID,
				AccessoryID: accessoryID,
				Exp:         0,
			})
			if err != nil {
				return err
			}
			existingCount[accessoryID]++
		}
	}

	return nil
}

func BuildAccessoryOwningUserIDMap(userSession *xorm.Session, userID int) (map[int][]int, error) {
	rows := []UserAccessory{}
	err := userSession.Table(new(UserAccessory)).
		Where("user_id = ?", userID).
		OrderBy("accessory_id ASC, accessory_owning_user_id ASC").
		Find(&rows)
	if err != nil {
		return nil, err
	}

	result := make(map[int][]int)
	for _, row := range rows {
		result[row.AccessoryID] = append(result[row.AccessoryID], row.AccessoryOwningUserID)
	}
	return result, nil
}

func ResolveAccessoryOwningUserIDsByAccessoryID(owningIDs map[int][]int, requiredCounts map[int]int) (map[int][]int, error) {
	result := make(map[int][]int, len(requiredCounts))
	for accessoryID, count := range requiredCounts {
		ids := owningIDs[accessoryID]
		if len(ids) < count {
			return nil, fmt.Errorf("insufficient user accessories for accessory_id=%d need=%d have=%d", accessoryID, count, len(ids))
		}
		result[accessoryID] = append([]int(nil), ids[:count]...)
	}

	for accessoryID := range result {
		sort.Ints(result[accessoryID])
	}

	return result, nil
}
