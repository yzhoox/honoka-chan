package secretbox

import (
	"fmt"
)

func SecretBoxApi(action string) (res any, err error) {
	switch action {
	case "all":
		res, err = all()
	default:
		err = fmt.Errorf("unimplemented action: secretbox: %s", action)
	}
	return res, err
}
