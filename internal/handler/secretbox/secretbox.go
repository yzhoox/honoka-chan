package secretbox

import (
	"fmt"
	unitmodel "honoka-chan/internal/model/unit"
	secretboxapischema "honoka-chan/internal/schema/api/secretbox"
	"honoka-chan/internal/session"
	honokautils "honoka-chan/internal/utils"
	"honoka-chan/pkg/db"
	"net/http"
	"sort"
	"time"
)

const secretBoxDataFile = "secretbox_data.json"

func loadSecretBoxAllData() (secretboxapischema.AllResp, error) {
	data, err := honokautils.LoadServerData[secretboxapischema.AllData](secretBoxDataFile)
	if err != nil {
		return secretboxapischema.AllResp{}, err
	}

	return secretboxapischema.AllResp{
		Result:     data,
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}, nil
}

func findSecretBoxPage(config secretboxapischema.AllResp, secretBoxID int) (int, secretboxapischema.SecretBoxPage, error) {
	for _, category := range config.Result.MemberCategoryList {
		for _, page := range category.PageList {
			if page.SecretBoxInfo.SecretBoxID == secretBoxID {
				return category.MemberCategory, page, nil
			}
		}
	}
	return 0, secretboxapischema.SecretBoxPage{}, fmt.Errorf("secretbox not found: %d", secretBoxID)
}

func findSecretBoxCost(page secretboxapischema.SecretBoxPage, costID int) (secretboxapischema.SecretBoxButton, secretboxapischema.SecretBoxCost, error) {
	for _, button := range page.ButtonList {
		for _, cost := range button.CostList {
			if cost.ID == costID {
				return button, cost, nil
			}
		}
	}
	return secretboxapischema.SecretBoxButton{}, secretboxapischema.SecretBoxCost{}, fmt.Errorf("secretbox cost not found: %d", costID)
}

func queryUnitIDs(ss *session.Session, memberCategory int, supportOnly bool, rarity int) ([]int, error) {
	var (
		rows []struct {
			UnitID int `xorm:"unit_id"`
		}
		err error
	)

	sql := `
SELECT u.unit_id
FROM unit_m u
LEFT JOIN unit_type_m ut ON ut.unit_type_id = u.unit_type_id
WHERE u.rarity = ?`
	args := []any{rarity}

	if supportOnly {
		sql += ` AND u.disable_rank_up != 0`
	} else if memberCategory >= 0 {
		sql += ` AND ut.member_category = ?`
		args = append(args, memberCategory)
	}

	if ss != nil {
		err = ss.MainEng.SQL(sql, args...).Find(&rows)
	} else {
		err = db.MainEng.SQL(sql, args...).Find(&rows)
	}
	if err != nil {
		return nil, err
	}

	result := make([]int, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.UnitID)
	}
	sort.Ints(result)
	return result, nil
}

func getUserUnitByUnitID(ss *session.Session, unitID int) (unitmodel.UnitDataMap, error) {
	unitData := unitmodel.UnitDataMap{}
	has, err := ss.GetBasicUnitInfo().
		Where("a.user_id = ?", ss.UserID).
		Where("a.unit_id = ?", unitID).
		Get(&unitData)
	if err != nil {
		return unitmodel.UnitDataMap{}, err
	}
	if !has {
		return unitmodel.UnitDataMap{}, fmt.Errorf("user %d does not own unit %d", ss.UserID, unitID)
	}
	return unitData, nil
}
