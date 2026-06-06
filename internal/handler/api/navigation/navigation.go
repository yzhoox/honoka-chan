package navigation

import (
	honokautils "honoka-chan/internal/utils"
)

func NavigationApi(action string) (res any, err error) {
	switch action {
	case "specialCutin":
		res, err = SpecialCutin()
	default:
		err = honokautils.NewUnimplementedActionError("navigation", action)
	}
	return res, err
}
