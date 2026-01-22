package costume

import "fmt"

func CostumeApi(action string) (res any, err error) {
	switch action {
	case "costumeList":
		res, err = costumeList()
	default:
		err = fmt.Errorf("unimplemented action: costume: %s", action)
	}
	return res, err
}
