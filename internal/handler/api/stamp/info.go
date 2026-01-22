package stamp

import (
	"encoding/json"
	honokautils "honoka-chan/pkg/utils"
)

func stampInfo() (res any, err error) {
	stampResp := honokautils.ReadAllText("assets/serverdata/stamp.json")
	err = json.Unmarshal([]byte(stampResp), &res)

	return res, err
}
