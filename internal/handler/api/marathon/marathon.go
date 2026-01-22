package marathon

import "fmt"

func MarathonApi(action string) (res any, err error) {
	switch action {
	case "marathonInfo":
		res, err = marathonInfo()
	default:
		err = fmt.Errorf("unimplemented action: marathon: %s", action)
	}
	return res, err
}
