package secretbox

import (
	honokautils "honoka-chan/internal/utils"
)

func SecretBoxApi(action string) (res any, err error) {
	switch action {
	case "all":
		res, err = all()
	default:
		err = honokautils.NewUnimplementedActionError("secretbox", action)
	}
	return res, err
}
