package banner

import honokautils "honoka-chan/internal/utils"

func BannerApi(action string) (res any, err error) {
	switch action {
	case "bannerList":
		res, err = bannerList()
	default:
		err = honokautils.NewUnimplementedActionError("banner", action)
	}
	return res, err
}
