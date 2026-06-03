package item

import (
	"encoding/json"
	"honoka-chan/pkg/utils"
)

func itemList() (res any, err error) {
	itemResp := utils.ReadAllText("assets/serverdata/item.json")
	err = json.Unmarshal([]byte(itemResp), &res)

	return res, err
}
