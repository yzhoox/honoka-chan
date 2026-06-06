package stamp

import honokautils "honoka-chan/internal/utils"

func StampApi(action string) (res any, err error) {
	switch action {
	case "stampInfo":
		res, err = stampInfo()
	default:
		err = honokautils.NewUnimplementedActionError("stamp", action)
	}
	return res, err
}
