package livese

import "fmt"

func LiveSeApi(action string) (res any, err error) {
	switch action {
	case "liveseInfo":
		res, err = LiveSeInfo()
	default:
		err = fmt.Errorf("unimplemented action: livese: %s", action)
	}
	return res, err
}
