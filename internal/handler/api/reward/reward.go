package reward

import (
	honokautils "honoka-chan/internal/utils"
)

func RewardApi(action string) (res any, err error) {
	switch action {
	case "rewardList":
		res, err = rewardList()
	case "rewardHistory":
		res, err = rewardHistory()
	default:
		err = honokautils.NewUnimplementedActionError("reward", action)
	}
	return res, err
}
