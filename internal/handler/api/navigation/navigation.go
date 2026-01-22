package navigation

import (
	"fmt"
)

func NavigationApi(action string) (res any, err error) {
	switch action {
	case "specialCutin":
		res, err = SpecialCutin()
	default:
		err = fmt.Errorf("unimplemented action: navigation: %s", action)
	}
	return res, err
}
