package item

import honokautils "honoka-chan/internal/utils"

func ItemApi(action string) (res any, err error) {
	switch action {
	case "list":
		res, err = itemList()
	default:
		err = honokautils.NewUnimplementedActionError("item", action)
	}
	return res, err
}
