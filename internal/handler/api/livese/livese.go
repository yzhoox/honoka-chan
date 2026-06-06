package livese

import honokautils "honoka-chan/internal/utils"

func LiveSeApi(action string) (res any, err error) {
	switch action {
	case "liveseInfo":
		res, err = LiveSeInfo()
	default:
		err = honokautils.NewUnimplementedActionError("livese", action)
	}
	return res, err
}
