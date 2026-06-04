package item

import (
	itemapischema "honoka-chan/internal/schema/api/item"
	honokautils "honoka-chan/internal/utils"
	"time"
)

func itemList() (res any, err error) {
	itemData, err := honokautils.LoadServerData[itemapischema.ListData]("item_data.json")
	if err != nil {
		return nil, err
	}

	res = itemapischema.ListResp{
		Result:     itemData,
		Status:     200,
		CommandNum: false,
		TimeStamp:  time.Now().Unix(),
	}

	return res, err
}
