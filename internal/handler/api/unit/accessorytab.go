package unit

import (
	"encoding/json"
	unitapischema "honoka-chan/internal/schema/api/unit"
	"honoka-chan/pkg/utils"
	"time"
)

func unitAccessoryTab() (res any, err error) {
	data := utils.ReadAllText("assets/serverdata/accessoryTab.json")
	resp := unitapischema.AccessoryTabResp{}
	err = json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return nil, err
	}
	resp.TimeStamp = time.Now().Unix()

	return resp, err
}
