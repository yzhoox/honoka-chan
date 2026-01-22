package liveicon

import "fmt"

func LiveIconApi(action string) (res any, err error) {
	switch action {
	case "liveiconInfo":
		res, err = liveIconInfo()
	default:
		err = fmt.Errorf("unimplemented action: liveicon: %s", action)
	}
	return res, err
}
