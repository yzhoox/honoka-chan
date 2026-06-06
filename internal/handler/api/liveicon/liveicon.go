package liveicon

import honokautils "honoka-chan/internal/utils"

func LiveIconApi(action string) (res any, err error) {
	switch action {
	case "liveiconInfo":
		res, err = liveIconInfo()
	default:
		err = honokautils.NewUnimplementedActionError("liveicon", action)
	}
	return res, err
}
