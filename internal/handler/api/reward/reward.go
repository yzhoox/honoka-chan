package reward

import (
	"fmt"
)

func RewardApi(action string) (res any, err error) {
	switch action {
	case "rewardList":
		res, err = rewardList()
	case "rewardHistory":
		res, err = rewardHistory()
	default:
		err = fmt.Errorf("unimplemented action: reward: %s", action)
	}
	return res, err
}
