package challenge

import (
	honokautils "honoka-chan/internal/utils"
)

func ChallengeApi(action string) (res any, err error) {
	switch action {
	case "challengeInfo":
		res, err = challengeInfo()
	default:
		err = honokautils.NewUnimplementedActionError("challenge", action)
	}
	return res, err
}
