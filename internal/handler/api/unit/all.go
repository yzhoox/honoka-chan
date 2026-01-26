package unit

import (
	unitmodel "honoka-chan/internal/model/unit"
	unitapischema "honoka-chan/internal/schema/api/unit"
	"honoka-chan/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

func unitAll(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	unitMap := []unitmodel.UnitDataMap{}
	err = ss.GetBasicUnitInfo().
		Where("a.user_id = ?", ss.UserID).Find(&unitMap)
	if ss.CheckErr(err) {
		return
	}

	unitsData := []unitapischema.Active{}
	for _, u := range unitMap {
		unitsData = append(unitsData, unitapischema.Active{
			UnitOwningUserID:            u.UnitOwningUserID,
			UserID:                      ss.UserID,
			UnitID:                      u.UnitID,
			Exp:                         u.Exp,
			NextExp:                     0,
			Level:                       u.Level,
			MaxLevel:                    u.MaxLevel,
			LevelLimitID:                u.LevelLimitID,
			Rank:                        u.Rank,
			MaxRank:                     u.MaxRank,
			Love:                        u.Love,
			MaxLove:                     u.MaxLove,
			UnitSkillExp:                u.UnitSkillExp,
			UnitSkillLevel:              u.UnitSkillLevel,
			MaxHp:                       u.MaxHp,
			UnitRemovableSkillCapacity:  u.UnitRemovableSkillCapacity,
			FavoriteFlag:                u.FavoriteFlag,
			DisplayRank:                 u.DisplayRank,
			IsRankMax:                   u.IsRankMax,
			IsLoveMax:                   u.IsLoveMax,
			IsLevelMax:                  u.IsLevelMax,
			IsSigned:                    u.IsSigned,
			IsSkillLevelMax:             u.IsSkillLevelMax,
			IsRemovableSkillCapacityMax: u.IsRemovableSkillCapacityMax,
			InsertDate:                  "2023-03-20 03:58:55",
		})
	}

	res = unitapischema.AllResp{
		Result: unitapischema.AllData{
			Active:  unitsData,
			Waiting: []unitapischema.Waiting{},
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
