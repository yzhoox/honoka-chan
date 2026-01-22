package item

import "fmt"

func ItemApi(action string) (res any, err error) {
	switch action {
	case "list":
		res, err = itemList()
	default:
		err = fmt.Errorf("unimplemented action: item: %s", action)
	}
	return res, err
}
