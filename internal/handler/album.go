package handler

import (
	"honoka-chan/internal/model"
	"honoka-chan/internal/session"
	"honoka-chan/internal/utils"

	"github.com/gin-gonic/gin"
)

type AlbumRes struct {
	UnitId int `xorm:"unit_id"`
	Rarity int `xorm:"rarity"`
}

func AlbumSeriesAll(ctx *gin.Context) {
	ss := session.New(ctx)
	defer ss.Finalize()

	var albumIds []int
	err := ss.MainEng.Table("album_series_m").Select("album_series_id").Find(&albumIds)
	if ss.CheckErr(err) {
		return
	}

	albumSeriesAllRes := []model.AlbumSeriesRes{}
	for _, albumId := range albumIds {
		unitList := []AlbumRes{}
		err = ss.MainEng.Table("unit_m").Where("album_series_id = ?", albumId).Cols("unit_id,rarity").Find(&unitList)
		if ss.CheckErr(err) {
			return
		}

		albumSeriesAll := []model.AlbumResult{}
		for _, unit := range unitList {
			albumSeries := model.AlbumResult{
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
				albumSeries.SignFlag = utils.IsSigned(unit.UnitId)
			}

			albumSeriesAll = append(albumSeriesAll, albumSeries)
		}

		albumSeriesAllRes = append(albumSeriesAllRes, model.AlbumSeriesRes{
			SeriesID: albumId,
			UnitList: albumSeriesAll,
		})
	}

	resp := model.AlbumSeriesResp{
		ResponseData: albumSeriesAllRes,
		ReleaseInfo:  []any{},
		StatusCode:   200,
	}

	ss.Respond(resp)
}
