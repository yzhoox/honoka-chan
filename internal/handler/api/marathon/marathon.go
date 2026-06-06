package marathon

import honokautils "honoka-chan/internal/utils"

func MarathonApi(action string) (res any, err error) {
	switch action {
	case "marathonInfo":
		res, err = marathonInfo()
	default:
		err = honokautils.NewUnimplementedActionError("marathon", action)
	}
	return res, err
}
