package banner

import "fmt"

func BannerApi(action string) (res any, err error) {
	switch action {
	case "bannerList":
		res, err = bannerList()
	default:
		err = fmt.Errorf("unimplemented action: banner: %s", action)
	}
	return res, err
}
