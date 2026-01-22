package album

import (
	"honoka-chan/internal/middleware"
	"honoka-chan/internal/router"
	"honoka-chan/internal/schema/album"
	"honoka-chan/internal/session"

	"github.com/gin-gonic/gin"
)

func seriesAll(ctx *gin.Context) {
	ss := session.Get(ctx)
	defer ss.Finalize()

	var albumID []int
	err := ss.MainEng.Table("album_series_m").Select("album_series_id").Find(&albumID)
	if ss.CheckErr(err) {
		return
	}

	albumSeriesAllRes := []album.SeriesAllData{}
	for _, id := range albumID {
		var unitList []struct {
			UnitId int `xorm:"unit_id"`
			Rarity int `xorm:"rarity"`
		}
		err = ss.MainEng.Table("unit_m").Where("album_series_id = ?", id).Cols("unit_id,rarity").Find(&unitList)
		if ss.CheckErr(err) {
			return
		}

		albumSeriesAll := []album.UnitList{}
		for _, unit := range unitList {
			albumSeries := album.UnitList{
				UnitID:           unit.UnitId,
				RankMaxFlag:      true,
				LoveMaxFlag:      true,
				RankLevelMaxFlag: true,
				AllMaxFlag:       true,
				TotalLove:        10000,
				FavoritePoint:    1000,
			}

			if unit.Rarity != 4 {
				switch unit.Rarity {
				case 1:
					// N
					albumSeries.HighestLovePerUnit = 50
				case 2:
					// R
					albumSeries.HighestLovePerUnit = 200
				case 3:
					// SR
					albumSeries.HighestLovePerUnit = 500
				case 5:
					// SSR
					albumSeries.HighestLovePerUnit = 750
				}
			} else {
				// UR
				albumSeries.HighestLovePerUnit = 1000

				// IsSigned
				albumSeries.SignFlag, err = ss.MainEng.Table("unit_sign_asset_m").Where("unit_id = ?", unit.UnitId).Exist()
				if ss.CheckErr(err) {
					return
				}

			}

			albumSeriesAll = append(albumSeriesAll, albumSeries)
		}

		albumSeriesAllRes = append(albumSeriesAllRes, album.SeriesAllData{
			SeriesID: id,
			UnitList: albumSeriesAll,
		})
	}

	ss.Respond(album.SeriesAllResp{
		ResponseData: albumSeriesAllRes,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	})
}

func init() {
	router.AddHandler("main.php", "POST", "/album/seriesAll", middleware.Common, seriesAll)
}
