package stamp

import (
	"encoding/json"
	"honoka-chan/pkg/utils"
)

func stampInfo() (res any, err error) {
	stampResp := utils.ReadAllText("assets/serverdata/stamp.json")
	err = json.Unmarshal([]byte(stampResp), &res)

	return res, err
}
