package album

import (
	albumapischema "honoka-chan/internal/schema/api/album"
	"honoka-chan/internal/session"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func albumAll(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	albumLists := []albumapischema.UnitList{}
	var unitList []struct {
		UnitId int `xorm:"unit_id"`
		Rarity int `xorm:"rarity"`
	}
	err = ss.MainEng.Table("unit_m").Cols("unit_id,rarity").OrderBy("unit_id ASC").Find(&unitList)
	if err != nil {
		return nil, err
	}

	for _, unit := range unitList {
		albumList := albumapischema.UnitList{
			RankMaxFlag:      true,
			LoveMaxFlag:      true,
			RankLevelMaxFlag: true,
			AllMaxFlag:       true,
			FavoritePoint:    1000,
		}
		albumList.UnitID = unit.UnitId
		if unit.Rarity != 4 {
			albumList.SignFlag = false
			switch unit.Rarity {
			case 1:
				albumList.HighestLovePerUnit = 50
				albumList.TotalLove = 50
			case 2:
				albumList.HighestLovePerUnit = 200
				albumList.TotalLove = 200
			case 3:
				albumList.HighestLovePerUnit = 500
				albumList.TotalLove = 500
			case 5:
				albumList.HighestLovePerUnit = 750
				albumList.TotalLove = 750
			}
		} else {
			albumList.HighestLovePerUnit = 1000
			albumList.TotalLove = 1000

			// IsSigned
			albumList.SignFlag, err = ss.MainEng.Table("unit_sign_asset_m").Where("unit_id = ?", unit.UnitId).Exist()
			if err != nil {
				return nil, err
			}
		}
		albumLists = append(albumLists, albumList)
	}

	res = albumapischema.AllResp{
		Result:     albumLists,
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
