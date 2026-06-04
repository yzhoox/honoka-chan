package unit

import (
	unitapischema "honoka-chan/internal/schema/api/unit"
	honokautils "honoka-chan/internal/utils"
	"time"
)

func unitAccessoryTab() (res any, err error) {
	tabList, err := honokautils.LoadServerData[[]unitapischema.TabList]("accessory_tab_list.json")
	if err != nil {
		return nil, err
	}

	res = unitapischema.AccessoryTabResp{
		Result: unitapischema.AccessoryTabData{
			TabList: tabList,
		},
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
