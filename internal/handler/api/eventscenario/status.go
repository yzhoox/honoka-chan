package eventscenario

import (
	"fmt"
	eventscenarioapischema "honoka-chan/internal/schema/api/eventscenario"
	"honoka-chan/internal/session"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func eventScenarioStatus(ctx *gin.Context) (res any, err error) {
	ss := session.Get(ctx)

	var eventID []int
	eventsList := []eventscenarioapischema.EventScenarioList{}
	err = ss.MainEng.Table("event_scenario_m").Cols("event_id").GroupBy("event_id").OrderBy("event_id DESC").Find(&eventID)
	if err != nil {
		return nil, err
	}

	for _, id := range eventID {
		var eventRes []struct {
			EventScenarioId int    `xorm:"event_scenario_id"`
			Chapter         int    `xorm:"chapter"`
			ChapterAsset    string `xorm:"chapter_asset"`
			OpenDate        string `xorm:"open_date"`
		}
		chapsList := []eventscenarioapischema.ChapterList{}
		err = ss.MainEng.Table("event_scenario_m").Where("event_id = ?", id).Cols("event_scenario_id,chapter,chapter_asset,open_date").
			OrderBy("chapter DESC").Find(&eventRes)
		if err != nil {
			return nil, err
		}

		for _, res := range eventRes {
			chapList := eventscenarioapischema.ChapterList{
				EventScenarioID: res.EventScenarioId,
				Chapter:         res.Chapter,
				ChapterAsset:    res.ChapterAsset,
				Status:          2,
				OpenFlashFlag:   0,
				IsReward:        false,
				CostType:        1000,
				ItemID:          1200,
				Amount:          1,
			}
			chapsList = append(chapsList, chapList)
		}

		event := eventscenarioapischema.EventScenarioList{
			EventID:     id,
			OpenDate:    strings.ReplaceAll(eventRes[0].OpenDate, "/", "-"),
			ChapterList: chapsList,
		}

		// HACK event_scenario_btn_asset
		switch id {
		case 10001:
			event.EventScenarioBtnAsset = "assets/image/ui/eventscenario/38_se_ba_t.png"
		case 221:
			event.EventScenarioBtnAsset = "assets/image/ui/eventscenario/215_se_ba_t.png"
		default:
			event.EventScenarioBtnAsset = fmt.Sprintf("assets/image/ui/eventscenario/%d_se_ba_t.png", id)
		}

		eventsList = append(eventsList, event)
	}
	res = eventscenarioapischema.StatusResp{
		Result: eventscenarioapischema.StatusData{
			EventScenarioList: eventsList,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
