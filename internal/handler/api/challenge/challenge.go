package challenge

import (
	"fmt"
)

func ChallengeApi(action string) (res any, err error) {
	switch action {
	case "challengeInfo":
		res, err = challengeInfo()
	default:
		err = fmt.Errorf("unimplemented action: challenge: %s", action)
	}
	return res, err
}
