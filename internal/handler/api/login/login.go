package login

import (
	"fmt"
)

func LoginApi(action string) (res any, err error) {
	switch action {
	case "topInfo":
		res, err = loginTopInfo()
	case "topInfoOnce":
		res, err = loginTopInfoOnce()
	default:
		err = fmt.Errorf("unimplemented action: login: %s", action)
	}
	return res, err
}
