package item

import (
	"encoding/json"
	honokautils "honoka-chan/pkg/utils"
)

func itemList() (res any, err error) {
	itemResp := honokautils.ReadAllText("assets/serverdata/item.json")
	err = json.Unmarshal([]byte(itemResp), &res)

	return res, err
}
