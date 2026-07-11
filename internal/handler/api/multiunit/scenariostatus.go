package multiunit

import (
	multiunitapischema "honoka-chan/internal/schema/api/multiunit"
	"honoka-chan/internal/session"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func MultiUnitScenarioStatus(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var statusID []int
	multiUnitsList := []multiunitapischema.StatusList{}
	err = ss.MainEng.Table("multi_unit_scenario_m").Cols("multi_unit_id").GroupBy("multi_unit_id").OrderBy("multi_unit_id ASC").Find(&statusID)
	if err != nil {
		return nil, err
	}

	for _, id := range statusID {
		var multiRes struct {
			MultiUnitScenarioId       int    `xorm:"multi_unit_scenario_id"`
			Chapter                   int    `xorm:"chapter"`
			MultiUnitScenarioBtnAsset string `xorm:"multi_unit_scenario_btn_asset"`
			OpenDate                  string `xorm:"open_date"`
		}
		_, err = ss.MainEng.Table("multi_unit_scenario_m").
			Join("LEFT", "multi_unit_scenario_open_m", "multi_unit_scenario_m.multi_unit_id = multi_unit_scenario_open_m.multi_unit_id").
			Cols("multi_unit_scenario_btn_asset,open_date,multi_unit_scenario_id,chapter").
			Where("multi_unit_scenario_m.multi_unit_id = ?", id).Get(&multiRes)
		if err != nil {
			return nil, err
		}

		multiUnitsList = append(multiUnitsList, multiunitapischema.StatusList{
			MultiUnitID:               id,
			Status:                    2,
			MultiUnitScenarioBtnAsset: multiRes.MultiUnitScenarioBtnAsset,
			OpenDate:                  strings.ReplaceAll(multiRes.OpenDate, "/", "-"),
			ChapterList: []multiunitapischema.ChapterList{
				{
					MultiUnitScenarioID: multiRes.MultiUnitScenarioId,
					Chapter:             multiRes.Chapter,
					Status:              2,
				},
			},
		})
	}
	res = multiunitapischema.StatusResp{
		Result: multiunitapischema.StatusData{
			MultiUnitScenarioStatusList:  multiUnitsList,
			UnlockedMultiUnitScenarioIds: []any{},
		},
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
