package stamp

import "fmt"

func StampApi(action string) (res any, err error) {
	switch action {
	case "stampInfo":
		res, err = stampInfo()
	default:
		err = fmt.Errorf("unimplemented action: stamp: %s", action)
	}
	return res, err
}
