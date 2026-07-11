package secretbox

import (
	secretboxapischema "honoka-chan/internal/schema/api/secretbox"
	honokautils "honoka-chan/internal/utils"
	"net/http"
	"time"
)

func all() (res any, err error) {
	secretBoxData, err := honokautils.LoadServerData[secretboxapischema.AllData]("secretbox_data.json")
	if err != nil {
		return nil, err
	}

	secretBoxData.UseCache = 0
	secretBoxData.IsUnitMax = false
	secretBoxData.GaugeInfo.MaxGaugePoint = 100
	secretBoxData.GaugeInfo.GaugePoint = 0

	for i := range secretBoxData.ItemList {
		if secretBoxData.ItemList[i].Amount < 999 {
			secretBoxData.ItemList[i].Amount = 999
		}
	}

	for i := range secretBoxData.MemberCategoryList {
		for j := range secretBoxData.MemberCategoryList[i].PageList {
			for k := range secretBoxData.MemberCategoryList[i].PageList[j].ButtonList {
				for m := range secretBoxData.MemberCategoryList[i].PageList[j].ButtonList[k].CostList {
					secretBoxData.MemberCategoryList[i].PageList[j].ButtonList[k].CostList[m].Payable = true
				}
			}
		}
	}

	res = secretboxapischema.AllResp{
		Result:     secretBoxData,
		Status:     http.StatusOK,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, nil
}
